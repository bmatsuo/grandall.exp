package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

func main() {
	configPath := flag.String("config", "", "configuration file")
	sitesDir := flag.String("sites", "", "sites directory")
	cssHRefs := flag.String("css", "", "css locations")
	flag.Parse()

	conf, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	sites, err := ReadSites(*sitesDir)
	if err != nil {
		log.Fatal(err)
	}

	hrefSep := func(c rune) bool { return strings.ContainsRune(", ", c) }
	css := strings.FieldsFunc(*cssHRefs, hrefSep)

	s := new(http.Server)
	s.Addr = conf.Bind
	s.Handler, err = RedirectHandler(sites, css)
	if err != nil {
		log.Fatal(err)
	}

	log.Panic(s.ListenAndServe())
}
