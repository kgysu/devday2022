package server

import (
	"github.com/gorilla/mux"
	"github.com/kgysu/devday2022/pkg/database"
	"net/http"
)

type MyServer struct {
	store *database.Directory
}

func NewServer(store *database.Directory) *MyServer {
	server := new(MyServer)
	server.store = store
	return server
}

func (s *MyServer) RegisterRoutes(router *mux.Router) {
	// Products
	productsRouter := router.PathPrefix("/products").Subrouter()
	productsRouter.HandleFunc("/{id:[a-zA-Z0-9|-]+}", s.productHandler).Methods(http.MethodGet, http.MethodDelete)
	productsRouter.HandleFunc("", s.productsHandler).Methods(http.MethodGet, http.MethodPost)
}
