package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var serviceURL = "https://packagecloud.io"
var token string

func init() {
	token = os.Getenv("PACKAGECLOUD_TOKEN")
}

func sendRequest(endpoint string, v interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", serviceURL+endpoint, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(token, "")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}

type distroVersion struct {
	DisplayName   string `json:"display_name"`
	ID            int    `json:"id"`
	IndexName     string `json:"index_name"`
	VersionNumber string `json:"version_number"`
}

type distro struct {
	DisplayName string          `json:"display_name"`
	IndexName   string          `json:"index_name"`
	Versions    []distroVersion `json:"versions"`
}

type distros map[string][]distro

func main() {
	var d distros
	err := sendRequest("/api/v1/distributions.json", &d)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", d)
}
