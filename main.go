package main

import (
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var SDE *sde.SDE

func init() {
	log.Println("Parsing templates")
	ParseTemplates()
}

func main() {
	log.Println("Loading SDE related things")
	var err error
	SDE, err = sde.Load("dust.sde")
	if err != nil {
		log.Fatalf("Couldn't open SDE file:( %v", err.Error())
	}
	m := mux.NewRouter()
	m.HandleFunc("/", HandlerIndex)
	m.HandleFunc("/info", HandlerInfo)
	m.HandleFunc("/search", HandlerSearch)
	m.HandleFunc("/type/{TypeID:[0-9]+}", HandlerType)
	m.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	log.Println("Starting http server.")
	http.ListenAndServe("0.0.0.0:1339", m)
}
