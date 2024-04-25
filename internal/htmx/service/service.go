package service

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strings"

	log "github.com/Minerva-System/minerva-go/pkg/log"
)

type Json map[string]interface{}

const BASEURL = "http://100.83.101.113:30000/api/v1"

func makeUrl(endpoint string) string {
	return fmt.Sprintf("%s%s", BASEURL, endpoint)
}

func makeUrlId(endpoint string, id string) string {
	return fmt.Sprintf("%s%s/%s", BASEURL, endpoint, id)
}

func GetCompanies() (data []Json, err error) {
	url := makeUrl("/companies")
	log.Info("GET %s", url)
	response, err := http.Get(url)
	if err != nil {
		log.Error("Error on request: %v", err)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	defer response.Body.Close()
	return
}

func NewCompany(form string) (data Json, err error) {
	url := makeUrl("/companies")
	log.Info("POST %s", url)
	response, err := http.Post(url, "application/json", strings.NewReader(form))
	if err != nil {
		log.Error("Error on request: %v", err)
		return
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	defer response.Body.Close()
	return
}

func DeleteCompany(id string) (err error) {
	url := makeUrlId("/companies", id)
	log.Info("DELETE %s", url)

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return
	}

	response, err := client.Do(req)
	if err != nil {
		log.Error("Error on request: %v", err)
		return
	}
	defer response.Body.Close()
	return
}
