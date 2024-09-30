package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	whatsMyIPv4URL = "http://ipv4.whatismyip.akamai.com"
	whatsMyIPv6URL = "http://ipv6.whatismyip.akamai.com"
	defaultTTL     = 300
)

type recordSet struct {
	Type   string   `json:"rrset_type"`
	TTL    int      `json:"rrset_ttl"`
	Values []string `json:"rrset_values"`
}

func main() {
	var apiKey, myIPv4, myIPv6, domainName, recordName string
	var ttl int

	flag.StringVar(&apiKey, "api_key", "", "API key for Gandi LiveDNS")
	flag.StringVar(&myIPv4, "ip4", getMyIPv4(), "IPv4 address to set. Default is to get your public IP from Akamai")
	flag.StringVar(&myIPv6, "ip6", getMyIPv6(), "IPv6 address to set. Default is to get your public IP from Akamai")
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

	if myIPv4 != "" {
		recordSet := recordSet{
			Type:   "A",
			TTL:    ttl,
			Values: []string{myIPv4},
		}
		updateRecord(apiKey, domainName, recordName, recordSet)
	} else {
		log.Print("No IPv4 was found or set")
	}

	if myIPv6 != "" {
		recordSet := recordSet{
			Type:   "AAAA",
			TTL:    ttl,
			Values: []string{myIPv6},
		}
		updateRecord(apiKey, domainName, recordName, recordSet)
	} else {
		log.Print("No IPv6 was found or set")
	}
}

func getMyIP(whatsMyIPURL string) string {
	resp, err := http.Get(whatsMyIPURL)
	if err != nil {
		log.Print(err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(body)
}

func getMyIPv4() string {
	return getMyIP(whatsMyIPv4URL)
}

func getMyIPv6() string {
	return getMyIP(whatsMyIPv6URL)
}

func updateRecord(apiKey string, domainName string, recordName string, recordSet recordSet) {
	log.Printf("Updating %s.%s to %+v", recordName, domainName, recordSet)
	url := fmt.Sprintf("https://api.gandi.net/v5/livedns/domains/%s/records/%s/%s", domainName, recordName, recordSet.Type)

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
	req.Header.Add("Authorization", fmt.Sprintf("Apikey %s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", string(responseBody))
}
