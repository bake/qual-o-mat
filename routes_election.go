package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"github.com/bake/qual-o-mat/qualomat"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func (s *server) handleElection() http.HandlerFunc {
	tmpl, err := template.ParseFiles("templates/main.html", "templates/election.html")
	if err != nil {
		return s.handleError(errors.Wrap(err, "could not parse template"), 500)
	}
	type response struct {
		Election   *qualomat.Election
		Statements []qualomat.Statement
		Answers    []qualomat.Answer
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "election id is not a number", 400)
			return
		}
		election, err := s.qom.Election(id)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		stmts, err := election.Statements()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		answers, err := election.Answers()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		tmpl.Execute(w, response{election, stmts, answers})
	}
}

type party struct {
	qualomat.Party
	Overlap int
}

type sortByOverlap []party

func (p sortByOverlap) Len() int           { return len(p) }
func (p sortByOverlap) Less(i, j int) bool { return p[i].Overlap < p[j].Overlap }
func (p sortByOverlap) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (s *server) handleElectionVote() http.HandlerFunc {
	tmpl, err := template.ParseFiles("templates/main.html", "templates/election_vote.html")
	if err != nil {
		return s.handleError(errors.Wrap(err, "could not parse template"), 500)
	}

	type response struct {
		Election   *qualomat.Election
		Parties    []party
		Statements []qualomat.Statement
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, "election id is not a number", 400)
			return
		}
		election, err := s.qom.Election(id)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get election: %v", err), 500)
			return
		}
		opinions, err := election.Opinions()
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get opinions: %v", err), 500)
			return
		}
		parties, err := election.Parties()
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get parties: %v", err), 500)
			return
		}
		partyMap := map[int]*party{}
		for _, p := range parties {
			partyMap[p.ID] = &party{Party: p}
		}
		statements, err := election.Statements()
		if err != nil {
			http.Error(w, fmt.Sprintf("could not get statements: %v", err), 500)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, fmt.Sprintf("could not parse form: %v", err), 400)
			return
		}
		for stmts, answers := range r.Form {
			if len(answers) == 0 {
				continue
			}
			stmt, err := strconv.Atoi(stmts)
			if err != nil {
				continue
			}
			answer, err := strconv.Atoi(answers[0])
			if err != nil {
				continue
			}
			for _, o := range opinions {
				if o.Statement == stmt && o.Answer == answer {
					partyMap[o.Party].Overlap++
				}
			}
		}

		res := response{election, []party{}, statements}
		for _, p := range partyMap {
			res.Parties = append(res.Parties, *p)
		}
		sort.Sort(sort.Reverse(sortByOverlap(res.Parties)))

		tmpl.Execute(w, res)
	}
}
