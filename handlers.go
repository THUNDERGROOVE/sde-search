package main

import (
	"errors"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/gorilla/mux"
	"html"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	Version string
	Branch  string
)

var HomeDirRegexp = regexp.MustCompile("/home/[a-zA-Z0-9]*/")

// Global is a struct given to every single Render call.
//
// It contains everything a page will hopefully ever need to render.  Makes
// templates way more consistant and easier to write.
type Global struct {
	SDECount   int
	SDEVersion string
	SDEOffical bool
	Types      []*sde.SDEType
	Type       *WrappedSDEType
	Devel      bool
	Error      error
	StackTrace template.HTML
	Version    string
	Branch     string
}

// NewGlobal returns a Global pointer and fills in some values that will be the same for every page.
func NewGlobal() *Global {
	return &Global{
		Devel:   Dev,
		Version: Version,
		Branch:  Branch,
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
			PassError(rw, req, err)
			return
		}
		Render(rw, "search.tmpl", g)
	case "GET":
		Render(rw, "search.tmpl", g)
	default:
		PassError(rw, req, errors.New("Invalid request method: "+req.Method))
		return
	}
}

func HandlerType(rw http.ResponseWriter, req *http.Request) {
	g := NewGlobal()

	if tids, ok := mux.Vars(req)["TypeID"]; ok {
		i, _ := strconv.Atoi(tids) // Ignore error because mux will ensure that it's castable to an int before letting the handdler kick in
		t, err := SDE.GetType(i)
		if err != nil {
			PassError(rw, req, err)
			return
		}

		g.Type = SDETypeToWraped(t)
		Render(rw, "type.tmpl", g)
	} else {
		PassError(rw, req, fmt.Errorf("TypeID %v is invalid", tids))
		return
	}

}

func HandlerTestPassError(rw http.ResponseWriter, req *http.Request) {
	err := errors.New("Test error :(")
	PassError(rw, req, err)
	return
}

// PassError renders an error page if the err is not nil
func PassError(rw http.ResponseWriter, req *http.Request, err error) {
	if err == nil {
		return
	}
	g := NewGlobal()
	g.Error = err
	buf := make([]byte, 1<<16)
	i := runtime.Stack(buf, true)
	// Escape string
	s := html.EscapeString(string(buf[:i]))
	// Replace line breaks with br
	s = strings.Replace(s, "\n", "<br>", -1)
	s = HomeDirRegexp.ReplaceAllString(s, "")
	g.StackTrace = template.HTML(s)
	Render(rw, "error.tmpl", g)
	rw.WriteHeader(http.StatusInternalServerError)
}

/*
	Development handlers
*/
func HandlerReload(rw http.ResponseWriter, req *http.Request) {
	if Dev {
		Templates = make(map[string]*template.Template)
		ParseTemplates()
	} else {
		PassError(rw, req, fmt.Errorf("Template reloading is only available in development mode."))
	}
}
