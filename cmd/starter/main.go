package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
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

	// Create a multiplex Router
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/products", ProductsHandler)

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

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(`[{"name":"product A"},{"name":"product B"}]`))
}
