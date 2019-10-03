package powerdns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type service struct {
	client *Client
}

// Client configuration structure
type Client struct {
	Scheme   string
	Hostname string
	Port     string
	VHost    string
	Headers  map[string]string

	httpClient *http.Client

	common service // Reuse a single struct instead of allocating one for each service on the heap

	Cryptokeys *CryptokeysService
	Records    *RecordsService
	Servers    *ServersService
	Statistics *StatisticsService
	Zones      *ZonesService
}

// logFatalf makes log.Fatalf testable
var logFatalf = log.Fatalf

// NewClient initializes a new client instance
func NewClient(baseURL string, vhost string, headers map[string]string, httpClient *http.Client) *Client {
	scheme, hostname, port, err := parseBaseURL(baseURL)
	if err != nil {
		logFatalf("%s is not a valid url: %v", baseURL, err)
	}

	c := &Client{
		Scheme:   scheme,
		Hostname: hostname,
		Port:     port,
		VHost:    parseVhost(vhost),
		Headers:  headers,
	}

	if httpClient == nil {
		c.httpClient = http.DefaultClient
	}

	c.common.client = c

	c.Cryptokeys = (*CryptokeysService)(&c.common)
	c.Records = (*RecordsService)(&c.common)
	c.Servers = (*ServersService)(&c.common)
	c.Statistics = (*StatisticsService)(&c.common)
	c.Zones = (*ZonesService)(&c.common)

	return c
}

func parseBaseURL(baseURL string) (string, string, string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", "", "", err
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

	return u.Scheme, hostname, port, nil
}

func parseVhost(vhost string) string {
	if vhost == "" {
		return "localhost"
	}
	return vhost
}

func generateAPIURL(scheme, hostname, port, path string) url.URL {
	u := url.URL{}
	u.Scheme = scheme
	u.Host = fmt.Sprintf("%s:%s", hostname, port)
	u.Path = fmt.Sprintf("/api/v1/%s", path)
	return u
}

func trimDomain(domain string) string {
	return strings.TrimSuffix(domain, ".")
}

func (p *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	apiURL := generateAPIURL(p.Scheme, p.Hostname, p.Port, path)
	req, err := http.NewRequest(method, apiURL.String(), buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	req.Header.Set("User-Agent", "go-powerdns")

	for key, value := range p.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}

func (p *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 401 {
		return resp, &Error{
			Status:  resp.Status,
			Message: "Unauthorized",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer func() {
			_ = resp.Body.Close()
		}()

		apiError := new(Error)
		err = json.NewDecoder(resp.Body).Decode(apiError)
		if err != nil {
			return resp, err
		}

		return resp, &Error{
			Status:  resp.Status,
			Message: apiError.Message,
		}
	}

	if v != nil {
		defer func() {
			_ = resp.Body.Close()
		}()

		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}

// Bool is a helper function that allocates a new bool value to store v and returns a pointer to it.
func Bool(v bool) *bool {
	return &v
}

// Uint32 is a helper function that allocates a new uint32 value to store v and returns a pointer to it.
func Uint32(v uint32) *uint32 {
	return &v
}

// Uint64 is a helper function that allocates a new uint64 value to store v and returns a pointer to it.
func Uint64(v uint64) *uint64 {
	return &v
}

// String is a helper function that allocates a new string value to store v and returns a pointer to it.
func String(v string) *string {
	return &v
}
