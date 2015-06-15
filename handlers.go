package main

import (
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/gorilla/mux"
	"html/template"
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
	Devel      bool
}

func NewGlobal() *Global {
	return &Global{
		Devel: Dev,
	}
}

func HandlerIndex(rw http.ResponseWriter, req *http.Request) {
	g := NewGlobal()
	Render(rw, "index.tmpl", g)
}

func HandlerInfo(rw http.ResponseWriter, req *http.Request) {
	g := NewGlobal()
	g.SDECount = len(SDE.Types)
	g.SDEVersion = SDE.Version
	g.SDEOffical = SDE.Official
	log.Println(g)
	Render(rw, "info.tmpl", g)
}

func HandlerSearch(rw http.ResponseWriter, req *http.Request) {
	g := NewGlobal()
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
	g := NewGlobal()

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

/*
	Development handlers
*/
func HandlerReload(rw http.ResponseWriter, req *http.Request) {
	if Dev {
		Templates = make(map[string]*template.Template)
		ParseTemplates()
	}
}
