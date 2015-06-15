package main

import (
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Global struct {
	SDECount   int
	SDEVersion string
	SDEOffical bool
	Types      []*sde.SDEType
	Type       *sde.SDEType
}

func HandlerIndex(rw http.ResponseWriter, req *http.Request) {
	g := &Global{}
	Render(rw, "index.tmpl", g)
}

func HandlerInfo(rw http.ResponseWriter, req *http.Request) {
	g := &Global{}
	g.SDECount = len(SDE.Types)
	g.SDEVersion = SDE.Version
	g.SDEOffical = SDE.Official
	log.Println(g)
	Render(rw, "info.tmpl", g)
}

func HandlerSearch(rw http.ResponseWriter, req *http.Request) {
	g := &Global{}
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			log.Println("Couldn't parse form", err.Error())
		}
		s := req.Form.Get("search")
		vals, err := SDE.Search(s)
		g.Types = vals
		if err != nil {
			log.Println("TODO: Cool error page.")
		}
		Render(rw, "search.tmpl", g)
	case "GET":
		Render(rw, "search.tmpl", g)
	default:
		log.Println("TODO: Make this notify the user somehow that this is wrong.")
	}
}

func HandlerType(rw http.ResponseWriter, req *http.Request) {
	g := &Global{}

	if tids, ok := mux.Vars(req)["TypeID"]; ok {
		i, _ := strconv.Atoi(tids) // Ignore error because mux will ensure that it's castable to an int before letting the handdler kick in
		t, err := SDE.GetType(i)
		if err != nil {
			log.Println("TODO: Show error page")
		}
		g.Type = t
		Render(rw, "type.tmpl", g)
	} else {
		log.Println("TODO: Error page when  TypeID doesn't exist")
	}

}
