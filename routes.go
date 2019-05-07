package main

import (
	"net/http"
)

func (s *server) routes() {
	s.router.HandleFunc("/election/{id:[0-9]+}", s.handleElection()).Methods(http.MethodGet)
	s.router.HandleFunc("/election/{id:[0-9]+}", s.handleElectionVote()).Methods(http.MethodPost)
	s.router.HandleFunc("/", s.handleElections())
	s.router.PathPrefix("/static").Handler(
		http.StripPrefix("/static", http.FileServer(http.Dir("./static/"))),
	)
}

func (s *server) handleError(err error, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.error(w, err, "could not initialize handler", code)
	}
}
