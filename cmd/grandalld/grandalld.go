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
	serviceHandler, err := service.Handler()
	if err != nil {
		log.Fatal(err)
	}

	s := new(http.Server)
	s.Addr = conf.Bind
	s.Handler = serviceHandler
	log.Panic(s.ListenAndServe())
}
