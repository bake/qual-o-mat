package main

import (
	"html/template"
	"net/http"
	"sort"

	"github.com/bake/qual-o-mat/qualomat"
	"github.com/pkg/errors"
)

type sortByDate []*qualomat.Election

func (e sortByDate) Len() int           { return len(e) }
func (e sortByDate) Less(i, j int) bool { return e[i].Date.Unix() < e[j].Date.Unix() }
func (e sortByDate) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func (s *server) handleElections() http.HandlerFunc {
	tmpl, err := template.ParseFiles("templates/main.html", "templates/elections.html")
	if err != nil {
		return s.handleError(errors.Wrap(err, "could not parse template"), 500)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	elections, err := s.qom.Elections()
	if err != nil {
		return s.handleError(errors.Wrap(err, "could not parse elections"), 500)
	}
	sort.Sort(sort.Reverse(sortByDate(elections)))
	type response struct{ Elections []*qualomat.Election }
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, response{elections})
	}
}
