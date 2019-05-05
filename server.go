package main

import (
	"net/http"
	"sync"

	"github.com/bake/qual-o-mat/qualomat"
	"github.com/gorilla/mux"
)

type server struct {
	qom    *qualomat.QualOMat
	mu     sync.Mutex
	router *mux.Router
}

func newServer(qom *qualomat.QualOMat) *server {
	s := &server{qom: qom, router: mux.NewRouter()}
	s.routes()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
