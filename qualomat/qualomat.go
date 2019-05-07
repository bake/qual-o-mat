package qualomat

import (
	"encoding/json"
	"os"
	"path"

	"github.com/pkg/errors"
)

type QualOMat struct {
	path      string
	elections []*Election
}

func New(path string) *QualOMat {
	return &QualOMat{path: path}
}

type Election struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Date date   `json:"date"`
	Path string `json:"path"`

	answers    []Answer
	comments   []Comment
	opinions   []Opinion
	overview   Overview
	parties    []Party
	statements []Statement
}

func (qom *QualOMat) Elections() ([]*Election, error) {
	var err error
	if qom.elections == nil {
		err = decode(path.Join(qom.path, "election.json"), &qom.elections)
		for i, e := range qom.elections {
			qom.elections[i].Path = path.Join(qom.path, e.Path)
		}
	}
	return qom.elections, err
}

func (qom *QualOMat) Election(id int) (*Election, error) {
	elections, err := qom.Elections()
	if err != nil {
		return nil, err
	}
	for _, e := range elections {
		if e.ID == id {
			return e, nil
		}
	}
	return nil, errors.Errorf("election %d not found", id)
}

func decode(file string, v interface{}) error {
	r, err := os.Open(file)
	if err != nil {
		return errors.Wrapf(err, "could not open %s", file)
	}
	defer r.Close()
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return errors.Wrapf(err, "could not decode %s", file)
	}
	return nil
}

type Answer struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func (e *Election) Answers() ([]Answer, error) {
	var err error
	if e.answers == nil {
		decode(path.Join(e.Path, "answer.json"), &e.answers)
	}
	return e.answers, err
}

type Comment struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

func (e *Election) Comments() ([]Comment, error) {
	var err error
	if e.comments == nil {
		err = decode(path.Join(e.Path, "comment.json"), &e.comments)
	}
	return e.comments, err
}

type Opinion struct {
	ID        int `json:"id"`
	Party     int `json:"party"`
	Statement int `json:"statement"`
	Answer    int `json:"answer"`
	Comment   int `json:"comment"`
}

func (e *Election) Opinions() ([]Opinion, error) {
	var err error
	if e.opinions == nil {
		err = decode(path.Join(e.Path, "opinion.json"), &e.opinions)
	}
	return e.opinions, err
}

type Overview struct {
	Title      string `json:"title"`
	Date       date   `json:"date"`
	Info       string `json:"info"`
	DataSource string `json:"data_source"`
}

func (e *Election) Overview() (Overview, error) {
	var err error
	if e.overview.Title == "" {
		err = decode(path.Join(e.Path, "overview.json"), &e.overview)
	}
	return e.overview, err
}

type Party struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Longname string `json:"longname"`
}

func (e *Election) Parties() ([]Party, error) {
	var err error
	if e.parties == nil {
		err = decode(path.Join(e.Path, "party.json"), &e.parties)
	}
	return e.parties, err
}

type Statement struct {
	ID       int    `json:"id"`
	Category int    `json:"category"`
	Label    string `json:"label"`
	Text     string `json:"text"`
}

func (e *Election) Statements() ([]Statement, error) {
	var err error
	if e.statements == nil {
		err = decode(path.Join(e.Path, "statement.json"), &e.statements)
	}
	return e.statements, err
}
