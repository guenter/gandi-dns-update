package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	whatsMyIPURL = "http://whatismyip.akamai.com"
	recordType   = "A"
	defaultTTL   = 300
)

type recordSet struct {
	Type   string   `json:"rrset_type"`
	TTL    int      `json:"rrset_ttl"`
	Values []string `json:"rrset_values"`
}

func main() {
	var apiKey, myIP, domainName, recordName string
	var ttl int

	flag.StringVar(&apiKey, "api_key", "", "API key for Gandi LiveDNS")
	flag.StringVar(&myIP, "ip", getMyIP(), "IP address to set. Default is to get your public IP from Akamai")
	flag.StringVar(&domainName, "domain_name", "", "Domain name")
	flag.StringVar(&recordName, "record_name", "", "Record name")
	flag.IntVar(&ttl, "ttl", defaultTTL, "TTL for the DNS record in seconds")
	flag.Parse()

	if apiKey == "" {
		log.Fatal("api_key can't be empty")
	}
	if domainName == "" {
		log.Fatal("domain_name can't be empty")
	}
	if recordName == "" {
		log.Fatal("record_name can't be empty")
	}

	recordSet := recordSet{
		Type:   recordType,
		TTL:    ttl,
		Values: []string{myIP},
	}

	updateRecord(apiKey, domainName, recordName, recordSet)
}

func getMyIP() string {
	resp, err := http.Get(whatsMyIPURL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func updateRecord(apiKey string, domainName string, recordName string, recordSet recordSet) {
	log.Printf("Updating %s.%s to %+v", recordName, domainName, recordSet)
	url := fmt.Sprintf("https://dns.api.gandi.net/api/v5/domains/%s/records/%s/%s", domainName, recordName, recordType)

	recordBytes, err := json.Marshal(recordSet)
	if err != nil {
		log.Fatal(err)
	}
	requestBody := bytes.NewReader(recordBytes)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, requestBody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", string(responseBody))
}
