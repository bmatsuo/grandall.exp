package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type API struct {
	mux     *http.ServeMux
	service *Service
}

func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.mux.ServeHTTP(w, r)
}

func APIHandler(service *Service) http.Handler {
	api := new(API)
	api.service = service
	api.mux = http.NewServeMux()
	api.mux.Handle("/v1/", http.StripPrefix("/v1/", APIHandlerV1(api.service)))
	return api
}

type APIV1 struct {
	mux     *http.ServeMux
	service *Service
}

func (api *APIV1) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.mux.ServeHTTP(w, r)
}

func APIHandlerV1(service *Service) http.Handler {
	v1 := new(API)
	v1.mux = http.NewServeMux()
	v1.mux.HandleFunc("/sites", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		switch r.Method {
		case "GET":
			err := json.NewEncoder(w).Encode(service.Sites())
			if err != nil {
				log.Printf("ERROR GET %s %v", r.URL.Path, err)
			}
		default:
			notAllowedV1("GET")(w, r)
		}
	})
	return v1
}

func notAllowedV1(allow string, other ...string) func(w http.ResponseWriter, r *http.Request) {
	all := append([]string{allow}, other...)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(all, " "))
		msg := fmt.Sprintf("%v request method must be %s", r.URL.Path, allow)
		http.Error(w, msg, http.StatusMethodNotAllowed)
	}
}
