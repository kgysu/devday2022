package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kgysu/devday2022/pkg/database"
	server2 "github.com/kgysu/devday2022/pkg/server"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

var LogFormatter = logrus.TextFormatter{
	DisableTimestamp: false,
	FullTimestamp:    true,
	TimestampFormat:  "2006-01-02 15:04:05 -07:00",
}

func main() {
	// Init logging
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&LogFormatter)

	// DB
	dir := createDir("postgresql://postgres:postgres@localhost:5432/devday")

	// Create a multiplex Router
	myServer := server2.NewServer(dir)
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler).Methods(http.MethodGet)
	router.Use(loggingMiddleware)
	myServer.RegisterRoutes(router)

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

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logrus.Infoln(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// Database connection

func createDir(dbUrlVal string) *database.Directory {
	var dbUrl pgURL
	err := dbUrl.Set(dbUrlVal)
	if err != nil {
		logrus.Fatalln("invalid db url:", err)
	}
	if dbUrl.String() == "" {
		logrus.Fatalln("no db url defined")
	}

	dir, err := database.NewDirectory(logrus.New(), (*url.URL)(&dbUrl))
	if err != nil {
		logrus.Fatalln("connection to db failed:", err)
	}

	return dir
}

type pgURL url.URL

func (p *pgURL) Set(in string) error {
	u, err := url.Parse(in)
	if err != nil {
		return err
	}

	switch u.Scheme {
	case "psql", "postgresql":
	default:
		return errors.New("unexpected scheme in URL")
	}

	*p = pgURL(*u)
	return nil
}

func (p pgURL) String() string {
	return (*url.URL)(&p).String()
}
