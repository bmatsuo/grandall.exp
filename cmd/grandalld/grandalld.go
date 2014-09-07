package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	influxdb "github.com/influxdb/influxdb/client"
)

func main() {
	configPath := flag.String("config", "", "configuration file")
	sitesDir := flag.String("sites", "", "sites directory")
	flag.Parse()

	conf, err := ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if conf.Metrics == "" {
		conf.Metrics = "log"
		conf.MetricsURL = "" // just to be tidy
	}

	var access func(name, bind, url string)
	switch conf.Metrics {
	case "log":
		access = LogAccess
	case "influxdb":
		client, err := InfluxDBClient(conf)
		if err != nil {
			log.Fatal(err)
		}
		access = InfluxDBAccess(client)
	default:
		log.Fatal("unknown metrics type")
	}

	sites, err := ReadSites(*sitesDir)
	if err != nil {
		log.Fatal(err)
	}

	s := new(http.Server)
	s.Addr = conf.Bind
	s.Handler, err = RedirectHandler(sites, access)
	if err != nil {
		log.Fatal(err)
	}

	log.Panic(s.ListenAndServe())
}

func LogAccess(name, bind, url string) { log.Printf("ACCESS %s %v", name, url) }

func InfluxDBAccess(c *influxdb.Client) func(name, bind, url string) {
	return func(name, bind, url string) {
		s := []*influxdb.Series{
			{
				Name:    "bind-access",
				Columns: []string{"name", "bind", "url"},
				Points: [][]interface{}{
					{name, bind, url},
				},
			},
		}
		err := c.WriteSeries(s)
		if err != nil {
			log.Printf("influxdb: %v", err)
		}
	}
}

// InfluxDBClient configures a new influxdb.Client using c.MetricsURL.
//	http://root:root@influxdb.example.com/mydb
//	https://root:root@influxdb.example.com/mydb
//	udp://root:root@influxdb.example.com/mydb
func InfluxDBClient(c *Config) (*influxdb.Client, error) {
	if c.Metrics != "influxdb" {
		return nil, fmt.Errorf("metrics are not reported to influxdb")
	}
	if c.MetricsURL == "" {
		return nil, fmt.Errorf("metrics url: blank")
	}
	u, err := url.Parse(c.MetricsURL)
	if err != nil {
		log.Fatalf("metrics url: %v", err)
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("metrics url: missing scheme")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("metrics url: invalid scheme")
	}

	conf := new(influxdb.ClientConfig)
	conf.Host = u.Host
	conf.Username = u.User.Username()
	conf.Password, _ = u.User.Password()
	conf.Database = strings.TrimPrefix(u.Path, "/")
	conf.IsSecure = u.Scheme == "https"
	conf.IsUDP = u.Scheme == "udp"

	return influxdb.NewClient(conf)
}
