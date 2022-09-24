package webserver

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jackc/pgtype"
	"github.com/kgysu/devday2022/pkg/database"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func (s *MyServer) productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.newProduct(w, r)
	case http.MethodGet:
		s.getProducts(w, r)
	}
}

func (s *MyServer) productHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		logrus.Errorln("no id specified")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		s.deleteProduct(w, id)
	case http.MethodGet:
		s.getProduct(w, id)
	}
}

func (s *MyServer) getProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	products, err := s.store.GetProducts(ctx)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (s *MyServer) newProduct(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var params database.AddProductParams
	err := decoder.Decode(&params)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	product, err := s.store.AddProduct(ctx, params)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (s *MyServer) getProduct(w http.ResponseWriter, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pgid pgtype.UUID
	err := pgid.Set(id)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	product, err := s.store.GetProduct(ctx, pgid)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (s *MyServer) deleteProduct(w http.ResponseWriter, id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pgid pgtype.UUID
	err := pgid.Set(id)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	product, err := s.store.DeleteProduct(ctx, pgid)
	if err != nil {
		logrus.Errorln(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}
