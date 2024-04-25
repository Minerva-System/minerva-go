package config

import (
	"fmt"
	"os"
	"strconv"
)

type MinervaHtmxConfig struct {
	Host       string
	Port       string
	FullHost   string
	ServerHost string
	Backend    string
	UseSSL     bool
}

var (
	Values MinervaHtmxConfig
)

func Load() {
	var ok bool
	var err error

	Values.Host, ok = os.LookupEnv("MINERVA_HTMX_HOST")
	if !ok {
		Values.Host = "127.0.0.1"
	}

	Values.Port, ok = os.LookupEnv("MINERVA_HTMX_PORT")
	if !ok {
		Values.Port = "5090"
	}

	use_ssl, ok := os.LookupEnv("MINERVA_HTMX_USE_SSL")
	if !ok {
		use_ssl = "false"
	}
	Values.UseSSL, err = strconv.ParseBool(use_ssl)
	if err != nil {
		Values.UseSSL = false
	}

	var protocol string = "http"
	if Values.UseSSL {
		protocol = "https"
	}

	Values.FullHost = fmt.Sprintf("%s://%s:%s", protocol, Values.Host, Values.Port)
	Values.ServerHost = fmt.Sprintf(":%s", Values.Port)

	Values.Backend, ok = os.LookupEnv("MINERVA_HTMX_BACKEND")
	if !ok {
		Values.Backend = "http://127.0.0.1:9000/api/v1"
	}

}
