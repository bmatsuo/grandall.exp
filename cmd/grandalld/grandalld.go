package main

import (
	"flag"
	"log"
	"net/http"

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
		conf.MetricsHost = "" // just to be tidy
	}

	var access func(name, bind, url string)
	switch conf.Metrics {
	case "log":
		access = LogAccess
	case "influxdb":
		if conf.MetricsHost == "" {
			log.Fatal("no influxdb host")
		}
		clientConf := &influxdb.ClientConfig{
			Host: conf.MetricsHost,
		}
		client, err := influxdb.NewClient(clientConf)
		if err != nil {
			log.Fatalf("influxdb: %v", err)
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
