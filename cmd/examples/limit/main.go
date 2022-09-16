package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

var limiter = NewRateLimiter(1, 1)

func main() {
	// Create a multiplex Router
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.Use(RateMiddleware)

	// Create HTTP Server
	addr := ":5000"
	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 2048,
	}
	logrus.Infoln("Server starting on ", addr)

	// Start listen
	logrus.Fatalln(server.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello DevDay 2022!")
}

func RateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// RateLimit per IP
		ip := r.RemoteAddr

		limiter := limiter.getVisitor(ip)
		if limiter.Allow() == false {
			logrus.Errorln("TooManyRequests from:", ip)
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
	limit    rate.Limit
	burst    int
}

// NewRateLimiter limit per second, burst: number to allow after limit exceeds within this second
func NewRateLimiter(limit rate.Limit, burst int) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		mu:       sync.Mutex{},
		limit:    limit,
		burst:    burst,
	}
}

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Run a background goroutine to remove old entries from the visitors map.
func (rl *RateLimiter) init() {
	go rl.cleanupVisitors()
}

func (rl *RateLimiter) getVisitor(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(rl.limit, rl.burst)
		// Include the current time when creating a new visitor.
		rl.visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}
