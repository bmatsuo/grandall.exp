/*
Command grandall is a client for grandalld.

Usage

The most common usage of grandall is to open a page in a web browser.  For
example if grandalld binds "/play" to "http://play.golang.org" the destination
can be opened with a grandall command.

	grandall play

If grandalld binds "/godoc" to "http://godoc.org" then both destinations can be
openend with one command.

	grandall play godoc

Configuration

Grandall looks for a TOML
configuration file pointing to grandalld.

	url = "http://localhost:7575/"

The default location of the configuration file is derived from the
environment, but a system wide fallback is used if no user/environment
specific configuration is found.
Grandall checks the following paths for configuration:

	$HOME/.config/grandall/grandall.conf
	/etc/grandall/grandall.conf

If no configuration file exists at one of the above locations then grandall
must be invoked with the "-config" flag.
*/
package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// DefaultConfigPath is a static path that is checked for a global
// configuration.  DefaultConfigPath is the lowest priority config and will not
// be read if a file is provided on the command line or in the user's home
// directory.
var DefaultConfigPath = "/etc/config"

func main() {
	configPath := flag.String("config", "", "TOML config file (otherwise $HOME/.config/grandall/grandall.conf or /etc/grandall/grandall.conf)")
	flag.Parse()

	config, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	baseURL, err := url.Parse(config.URL)
	if err != nil {
		log.Fatal(err)
	}
	if baseURL.Host == "" {
		log.Fatalf("unknown grandalld host: %v", baseURL)
	}
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	for _, urlstr := range flag.Args() {
		u, err := url.Parse(urlstr)
		if err != nil {
			log.Printf("%q %v", urlstr, err)
			continue
		}
		if u.Scheme == "" {
			u.Scheme = baseURL.Scheme
		}
		if u.Scheme == "" {
			u.Scheme = "http"
		}
		if u.Host == "" {
			u.Host = baseURL.Host
		}
		u.Path = baseURL.Path + strings.TrimPrefix(u.Path, "/")
		err = OpenURL(u.String())
		if err != nil {
			log.Printf("%q %v", urlstr, err)
			continue
		}
	}
}

func ReadConfig(filename string) (*Config, error) {
	if filename != "" {
		return readConfig(filename)
	}
	home := os.Getenv("HOME")
	if home != "" {
		filename = filepath.Join(home, ".config", "grandall", "grandall.conf")
	}
	c, err := readConfig(filename)
	if os.IsNotExist(err) {
		c, err = readConfig(DefaultConfigPath)
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func readConfig(filename string) (*Config, error) {
	c := new(Config)
	_, err := toml.DecodeFile(filename, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

type Config struct {
	URL string `toml:"url"`
}
