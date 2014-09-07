package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	configPath := flag.String("config", "", "configuration file")
	sitesDir := flag.String("sites", "", "sites directory")
	flag.Parse()

	conf, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	sites, err := ReadSites(*sitesDir)
	if err != nil {
		log.Fatal(err)
	}

	service, err := NewService(sites)
	if err != nil {
		log.Fatal(err)
	}
	service.Access = func(name, bind, url string) {
		log.Printf("ACCESS %s %v", name, url)
	}
	service.NotFound = func(w http.ResponseWriter, r *http.Request) {
		log.Printf("NOTFOUND %v", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	s := new(http.Server)
	s.Addr = conf.Bind
	s.Handler = service.Handler()
	log.Panic(s.ListenAndServe())
}
