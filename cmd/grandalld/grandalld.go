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

	s := new(http.Server)
	s.Addr = conf.Bind
	access := func(name, bind, url string) { log.Printf("ACCESS %s %v", name, url) }
	s.Handler, err = RedirectHandler(sites, access)
	if err != nil {
		log.Fatal(err)
	}

	log.Panic(s.ListenAndServe())
}
