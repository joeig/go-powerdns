package powerdns

import (
	"fmt"
	"github.com/dghubble/sling"
	"log"
	"net/url"
	"strings"
)

type Error struct {
	Message string `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

type PowerDNS struct {
	Scheme   string
	Hostname string
	Port     string
	VHost    string
	Domain   string
	APIKey   string
}

func NewClient(baseURL string, vhost string, domain string, apikey string) *PowerDNS {
	if vhost == "" {
		vhost = "localhost"
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("%s is not a valid url: %v", baseURL, err)
	}
	hp := strings.Split(u.Host, ":")
	hostname := hp[0]
	var port string
	if len(hp) > 1 {
		port = hp[1]
	} else {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	return &PowerDNS{
		Scheme:   u.Scheme,
		Hostname: hostname,
		Port:     port,
		VHost:    vhost,
		Domain:   domain,
		APIKey:   apikey,
	}
}

func (p *PowerDNS) makeSling() *sling.Sling {
	u := url.URL{}
	u.Host = p.Hostname + ":" + p.Port
	u.Scheme = p.Scheme
	u.Path = "/api/v1/"

	mySling := sling.New()
	mySling.Base(u.String())
	mySling.Set("X-API-Key", p.APIKey)

	return mySling
}
