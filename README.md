# Gandi DNS Update Tool

Updates an A record via the [Gandi LiveDNS API](https://doc.livedns.gandi.net/), by default using your public IP on the Internet.
This is an easy alternative to dynamic DNS services.

## Install

```
go install
```

## Usage

```
$ gandi-dns-update -help
Usage of gandi-dns-update:
  -api_key string
    	API key for Gandi LiveDNS
  -domain_name string
    	Domain name
  -ip string
    	IP address to set. Default is to get your public IP from Akamai (default "YOUR_IP")
  -record_name string
    	Record name
  -ttl int
    	TTL for the DNS record in seconds (default 300)
```