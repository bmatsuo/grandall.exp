package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

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
	Name        string `toml:"-" json:"name"`
	Bind        string `toml:"bind" json:"bind"`
	URL         string `toml:"url" json:"url"`
	Description string `toml:"description" json:"description"`
}

func BindURL(s *Site) (*url.URL, error) {
	u, err := url.Parse(s.Bind)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(u.Path, "/") {
		return nil, fmt.Errorf("relative bind")
	}
	if strings.HasSuffix(u.Path, "/") {
		return nil, fmt.Errorf("rooted bind")
	}
	return u, nil
}
