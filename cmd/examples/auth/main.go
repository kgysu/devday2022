package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Example to defend against time-based attacks on basic authentication
func main() {
	// Create a multiplex Router
	router := mux.NewRouter()
	router.HandleFunc("/", ExampleHandler)

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

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	// username ignored for this example
	_, pw, ok := r.BasicAuth()
	if !ok {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Check if Password matches, timing safe implementation
	passwordHash := sha256.Sum256([]byte(pw))
	expectedPasswordHash := sha256.Sum256([]byte("some")) // load pw from db
	passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1
	if !passwordMatch {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK!")
}
