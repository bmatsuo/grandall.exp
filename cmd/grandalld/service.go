package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// Service exposes the HTTP redirection service and an API for managing the
// service and creating user interfaces.
type Service struct {
	Access   func(name, bind, url string)
	NotFound func(w http.ResponseWriter, r *http.Request)
	sites    []*Site
	index    http.Handler
}

// NewService validates sites and allocates a new Service aliasing them.
func NewService(sites []*Site) (*Service, error) {
	names := make(map[string]bool)
	binds := make(map[string]bool)

	for _, site := range sites {
		// validate destination and bind urls
		u, err := url.Parse(site.URL)
		if err != nil {
			return nil, fmt.Errorf("%q: %v", site.URL, err)
		}
		if u.Scheme == "" {
			return nil, fmt.Errorf("relative url")
		}
		u, err = BindURL(site)
		if err != nil {
			return nil, fmt.Errorf("%q: %v", site.Name, err)
		}

		if names[site.Name] {
			return nil, fmt.Errorf("%q: duplicate name", site.Name)
		}
		names[site.Name] = true

		if binds[site.Bind] {
			return nil, fmt.Errorf("%q: duplicate bind %q", site.Name, site.Bind)
		}
		binds[site.Bind] = true
	}

	s := &Service{
		sites: sites,
		index: UI(sites),
	}

	return s, nil
}

func (s *Service) Sites() []*Site {
	return append([]*Site(nil), s.sites...)
}

// Handler constructs and returns an http.Handler for the service.  The
// returned handler provides alias endpoints redirecting to destination URLs,
// an API, and a static user interface at "/".
func (s *Service) Handler() http.Handler {
	root := http.NewServeMux()
	root.Handle("/.api/", http.StripPrefix("/.api", APIHandler(s)))
	root.Handle("/static/", s.index)
	redir := s.redirectHandler()
	redir.NotFoundHandler = http.HandlerFunc(s.notFound)
	root.Handle("/", redir)
	return root
}

func (service *Service) redirectHandler() *mux.Router {
	rt := mux.NewRouter()
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
			if service.Access != nil {
				service.Access(site.Name, site.Bind, site.URL)
			}
			http.Redirect(w, r, site.URL, http.StatusTemporaryRedirect)
		})
	}
	return rt
}

func (service *Service) notFound(w http.ResponseWriter, r *http.Request) {
	if service.NotFound != nil {
		service.NotFound(w, r)
		return
	}
	http.NotFound(w, r)
}
