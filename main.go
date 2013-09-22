package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pearkes/sv-frontend/data"
	"github.com/pearkes/sv-frontend/stats"
	"log"
	"net/http"
	"os"
)

var host = os.Getenv("HOST")
var db *data.Orm = nil
var red *data.Red = nil
var metrics *stats.StatsSink = nil
var DROPBOX_KEY = os.Getenv("DROPBOX_KEY")
var DROPBOX_SECRET = os.Getenv("DROPBOX_SECRET")
var DROPBOX_CALLBACK = os.Getenv("DROPBOX_CALLBACK")

func main() {
	db = data.NewOrm(os.Getenv("DATABASE_CONNECTION"))
	red = data.NewRedis(os.Getenv("REDIS_ADDRESS"), os.Getenv("REDIS_AUTH"))
	metrics = stats.NewStatsSink(os.Getenv("LIBRATO_USER"), os.Getenv("LIBRATO_TOKEN"), stats.ENV_WEB)

	// Bootstrap the database if the flag was given
	bootstrap := flag.Bool("bootstrap", false, "bootstrap the database")
	flag.Parse()
	if *bootstrap == true {
		db.Create()
		log.Println("Bootstrapping database...")
		os.Exit(0) // done, exit
	}

	r := mux.NewRouter()

	// Create a subroute for only our host, for the smallvictori.es page
	s := r.Host(host).Subrouter()
	// The main homepage
	s.HandleFunc("/", homeHandler)
	// The "help and support" page
	s.HandleFunc("/help", helpHandler)
	// Recieves callbacks from Dropbox during OAuth
	s.HandleFunc("/dbx/auth", authCallbackHandler)
	// Initializes the Dropbox OAuth flow
	s.HandleFunc("/dbx/init", authInitHandler)

	// Catch-all user page handler, for subdomains
	r.HandleFunc("/", userPageHandler)
	http.Handle("/", r)

	log.Println("Listening for requests...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(fmt.Sprintf("Server failed to start: %s", err.Error()))
	}
}
