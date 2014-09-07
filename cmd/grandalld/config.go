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

	// Metrics specifies the type of metrics to use (e.g. "influxdb", "log",
	// etc), the default value is "log".
	Metrics string `toml:"metrics"`
	// MetricsURL is ignored for "log" metrics.
	MetricsURL string `toml:"metrics_url"`
}

type Site struct {
	Name        string `toml:"-"`
	Bind        string `toml:"bind"`
	URL         string `toml:"url"`
	Description string `toml:"description"`
}

func BindURL(s *Site) (*url.URL, error) {
	u, err := url.Parse(s.Bind)
	if err != nil {
		return nil, err
	}
	return u, nil
}
