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
	service := &Service{
		sites: sites,
	}

	redirector, err := RedirectHandler(service)
	if err != nil {
		log.Panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/.api/", http.StripPrefix("/.api/", APIHandler(service)))
	mux.Handle("/", redirector)
	s := new(http.Server)
	s.Addr = conf.Bind
	s.Handler = mux

	log.Panic(s.ListenAndServe())
}

type Service struct {
	sites []*Site
}

func (s *Service) Sites() []*Site {
	return append([]*Site(nil), s.sites...)
}
