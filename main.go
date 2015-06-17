package main

import (
	"github.com/THUNDERGROOVE/SDETool/sde"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var (
	SDE *sde.SDE
	Dev = false
)

func init() {
	if _, err := os.Stat(".git"); err == nil {
		log.Println("Git data found.  Running in development mode")
		Dev = true
	}
	if Dev == false {
		log.Printf("In non-development environment.  Unpacking assets")
		UnpackAssets()
	}
	log.Println("Parsing templates")
	ParseTemplates()
	log.Printf("SDE Search %v@%v", Version, Branch)
	log.Println("Loading SDE related things")
	var err error
	SDE, err = sde.Load("dust.sde")
	if err != nil {
		log.Fatalf("Couldn't open SDE file:( %v", err.Error())
	}
}

func main() {
	m := mux.NewRouter()
	m.HandleFunc("/", HandlerIndex)
	m.HandleFunc("/info", HandlerInfo)
	m.HandleFunc("/search", HandlerSearch)
	m.HandleFunc("/type/{TypeID:[0-9]+}", HandlerType)
	m.HandleFunc("/store", HandlerStoreView)
	m.HandleFunc("/error", HandlerTestPassError)

	// Devel stuff
	m.HandleFunc("/dev/reload", HandlerReload)
	m.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	log.Println("Starting http server.")
	http.ListenAndServe("0.0.0.0:1339", m)
}
