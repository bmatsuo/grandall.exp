/*
Command grandalld provides an HTTP aliasing service.

Usage

Full usage information can be found by invoking grandalld with the -h flag.

	grandalld -h

A typical invocation looks like

	grandalld -config=/etc/grandall/grandalld.conf -sites=/etc/grandall/sites-enabled

Flags can also be specified through environment variables.

	export GRANDALLD_CONFIG=/etc/grandall/grandalld.conf
	export GRANDALLD_SITES=/etc/grandall/sites-enabled
	grandalld

Configuration

Grandalld configuration is specified in TOML format file (grandall.conf).

	bind    The address grandalld binds to for HTTP service.

Sites

Site aliases are defined in TOML format files.

	bind         The url grandalld binds the alias to.
	url          The url requests to the bind url are redirected to.
	description  Optional. A brief description of the destination.

Alias definitions are located in a directory specified either by the flag
-sites or by the GRANDALLD_SITES environment variable.  Aliases must not be
nested under subdirectories of GRANDALLD_SITES.  Site aliases are referred to
internally by the site's unique basename.
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {
	configPath := flag.String("config", os.Getenv("GRANDALLD_CONFIG"), "configuration file")
	sitesDir := flag.String("sites", os.Getenv("GRANDALLD_SITES"), "sites directory")
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
