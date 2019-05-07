package main

import (
	"fmt"
	"html/template"
	"math"
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
		Overview   qualomat.Overview
		Statements []qualomat.Statement
		Answers    []qualomat.Answer
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var res response
		var err error

		s.mu.Lock()
		defer s.mu.Unlock()
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			s.error(w, errors.New("election id is not a number"), fmt.Sprintf("election id %v is not a number", mux.Vars(r)["id"]), 400)
			return
		}
		if res.Election, err = s.qom.Election(id); err != nil {
			s.error(w, err, "could not get election", 500)
			return
		}
		if res.Overview, err = res.Election.Overview(); err != nil {
			s.error(w, err, "could not get overview", 500)
			return
		}
		if res.Statements, err = res.Election.Statements(); err != nil {
			s.error(w, err, "could not get statements", 500)
			return
		}
		if res.Answers, err = res.Election.Answers(); err != nil {
			s.error(w, err, "could not get answers", 500)
			return
		}
		tmpl.Execute(w, res)
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
		Overview   qualomat.Overview
		Parties    []party
		Statements []qualomat.Statement
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			s.error(w, err, "election id is not a number", 400)
			return
		}
		election, err := s.qom.Election(id)
		if err != nil {
			s.error(w, err, "could not get election", 500)
			return
		}
		overview, err := election.Overview()
		if err != nil {
			s.error(w, err, "could not get overview", 500)
			return
		}
		answers, err := election.Answers()
		if err != nil {
			s.error(w, err, "could not get answers", 500)
			return
		}
		opinions, err := election.Opinions()
		if err != nil {
			s.error(w, err, "could not get opinions", 500)
			return
		}
		parties, err := election.Parties()
		if err != nil {
			s.error(w, err, "could not get parties", 500)
			return
		}
		partyMap := map[int]*party{}
		for _, p := range parties {
			partyMap[p.ID] = &party{Party: p}
		}
		statements, err := election.Statements()
		if err != nil {
			s.error(w, err, "could not get statements", 500)
			return
		}

		if err := r.ParseForm(); err != nil {
			s.error(w, err, "could not parse form", 400)
			return
		}
		for stmtVal, answerVal := range r.Form {
			if len(answers) == 0 {
				continue
			}
			stmt, err := strconv.Atoi(stmtVal)
			if err != nil {
				continue
			}
			answer, err := strconv.Atoi(answerVal[0])
			if err != nil {
				continue
			}
			if answer == len(answers) {
				continue
			}
			for _, o := range opinions {
				if o.Statement == stmt {
					partyMap[o.Party].Overlap += 2 - int(math.Abs(float64(o.Answer-answer)))
				}
			}
		}

		res := response{election, overview, []party{}, statements}
		for _, p := range partyMap {
			res.Parties = append(res.Parties, *p)
		}
		sort.Sort(sort.Reverse(sortByOverlap(res.Parties)))

		tmpl.Execute(w, res)
	}
}
