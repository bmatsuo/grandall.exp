package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

func RedirectHandler(service *Service) (http.Handler, error) {
	sites := service.Sites()
	index, err := Index(sites)
	if err != nil {
		return nil, err
	}

	router := mux.NewRouter()
	for _, site := range sites {
		site := site
		_, err := url.Parse(site.URL)
		if err != nil {
			return nil, fmt.Errorf("%q %v", site.URL, err)
		}
		bindurl, err := BindURL(site)
		if err != nil {
			return nil, fmt.Errorf("%q %v", site.Name, err)
		}
		r := router.NewRoute()
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
	router.Handle("/", index)
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body.Close()
		log.Printf("NOTFOUND %v", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
	return router, nil
}
