package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/bake/qual-o-mat/qualomat"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
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

func (s *server) error(w http.ResponseWriter, err error, message string, code int) {
	http.Error(w, message, code)
	if message != "" {
		err = errors.Wrap(err, message)
	}
	log.Println(err)
}
