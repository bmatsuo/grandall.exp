package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/bbangert/toml"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter()
	for _, site := range sites {
		_, err := url.Parse(site.URL)
		if err != nil {
			log.Fatalf("%q %v", site.URL, err)
		}
		bindurl, err := BindURL(site)
		if err != nil {
			log.Fatalf("%q %v", site.Name, err)
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
			log.Printf("ACCESS %s %v", site.Name, site.URL)
			http.Redirect(w, r, site.URL, http.StatusTemporaryRedirect)
		})
	}
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("NOTFOUND %v", r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})

	log.Panic(http.ListenAndServe(conf.Bind, router))
}

func ReadConfig(filename string) (*Config, error) {
	c := new(Config)
	_, err := toml.DecodeFile(filename, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func ReadSites(dir string) ([]*Site, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sites := make([]*Site, 0, len(files))
	for _, info := range files {
		site, err := readSite(filepath.Join(dir, info.Name()))
		if err == errNotSite {
			continue
		}
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}
	return sites, err
}

var errNotSite = fmt.Errorf("not a site")

func readSite(filename string) (*Site, error) {
	base := filepath.Base(filename)
	if base == "README" {
		return nil, errNotSite
	}
	if strings.HasPrefix(base, ".") {
		return nil, errNotSite
	}
	site := new(Site)
	_, err := toml.DecodeFile(filename, site)
	if err != nil {
		return nil, err
	}
	site.Name = filepath.Base(filename)
	return site, nil
}

type Config struct {
	Bind string `toml:"bind"`
}

type Site struct {
	Name string `toml:"-"`
	Bind string `toml:"bind"`
	URL  string `toml:"url"`
}

func BindURL(s *Site) (*url.URL, error) {
	u, err := url.Parse(s.Bind)
	if err != nil {
		return nil, err
	}
	return u, nil
}
