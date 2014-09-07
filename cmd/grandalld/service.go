package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

type Service struct {
	sites []*Site
	index http.Handler
}

func NewService(sites []*Site) (*Service, error) {
	var err error
	for _, site := range sites {
		// validate destination and bind urls
		_, err = url.Parse(site.URL)
		if err != nil {
			return nil, fmt.Errorf("%q %v", site.URL, err)
		}
		_, err = BindURL(site)
		if err != nil {
			return nil, fmt.Errorf("%q %v", site.Name, err)
		}
	}

	s := &Service{sites: sites}

	s.index, err = Index(sites)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Service) Sites() []*Site {
	return append([]*Site(nil), s.sites...)
}

func (s *Service) Handler() (http.Handler, error) {
	root := http.NewServeMux()
	root.Handle("/.api/", http.StripPrefix("/.api", APIHandler(s)))

	rt := mux.NewRouter()
	root.Handle("/", rt)
	s.bindRedirects(rt)
	rt.Handle("/", s.index)
	rt.NotFoundHandler = http.HandlerFunc(s.notFound)
	return root, nil
}

func (service *Service) bindRedirects(rt *mux.Router) {
	for _, site := range service.sites {
		site := site

		// error checked at the time of construction
		bindurl, _ := BindURL(site)
		r := rt.NewRoute()
		r.Methods("GET")
		if bindurl.Scheme != "" {
			r.Schemes(bindurl.Scheme)
		}
		if bindurl.Host != "" {
			r.Host(bindurl.Host)
		}
		if bindurl.Path != "" {
			r.Path(bindurl.Path)
		}
		r.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body.Close()
			log.Printf("ACCESS %s %v", site.Name, site.URL)
			http.Redirect(w, r, site.URL, http.StatusTemporaryRedirect)
		})
	}
}

func (service *Service) notFound(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
	log.Printf("NOTFOUND %v", r.URL.Path)
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
