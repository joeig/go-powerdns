package powerdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// NewOption is a functional option for New.
type NewOption func(*Client)

// WithHeaders is an option for New to set HTTP client headers.
func WithHeaders(headers map[string]string) NewOption {
	return func(client *Client) {
		client.Headers = headers
	}
}

// WithHTTPClient is an option for New to set an HTTP client.
func WithHTTPClient(httpClient *http.Client) NewOption {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

// WithAPIKeyHeader is an option for New to set the API key header.
func WithAPIKeyHeader(key string) NewOption {
	return func(client *Client) {
		if client.Headers == nil {
			client.Headers = map[string]string{}
		}

		client.Headers["X-API-Key"] = key
	}
}

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

	Config     *ConfigService
	Cryptokeys *CryptokeysService
	Records    *RecordsService
	Servers    *ServersService
	Statistics *StatisticsService
	Zones      *ZonesService
	TSIGKey    *TSIGKeyService
}

// logFatalf makes log.Fatalf testable
var logFatalf = log.Fatalf

// NewClient initializes a new client instance.
// Deprecated: Use New with functional options instead.
// NewClient will be removed with the next major version.
func NewClient(baseURL string, vHost string, headers map[string]string, httpClient *http.Client) *Client {
	selectedHttpClient := httpClient
	if httpClient == nil {
		selectedHttpClient = http.DefaultClient
	}

	return New(baseURL, vHost, WithHeaders(headers), WithHTTPClient(selectedHttpClient))
}

// New initializes a new client instance.
func New(baseURL string, vHost string, options ...NewOption) *Client {
	scheme, hostname, port, err := parseBaseURL(baseURL)
	if err != nil {
		logFatalf("%s is not a valid url: %v", baseURL, err)
	}

	client := &Client{
		Scheme:     scheme,
		Hostname:   hostname,
		Port:       port,
		VHost:      parseVHost(vHost),
		httpClient: http.DefaultClient,
	}

	client.common.client = client

	client.Config = (*ConfigService)(&client.common)
	client.Cryptokeys = (*CryptokeysService)(&client.common)
	client.Records = (*RecordsService)(&client.common)
	client.Servers = (*ServersService)(&client.common)
	client.Statistics = (*StatisticsService)(&client.common)
	client.Zones = (*ZonesService)(&client.common)
	client.TSIGKey = (*TSIGKeyService)(&client.common)

	for _, option := range options {
		option(client)
	}

	return client
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

func parseVHost(vHost string) string {
	if vHost == "" {
		return "localhost"
	}
	return vHost
}

func generateAPIURL(scheme, hostname, port, path string, query *url.Values) url.URL {
	u := url.URL{}
	u.Scheme = scheme
	u.Host = fmt.Sprintf("%s:%s", hostname, port)
	u.Path = fmt.Sprintf("/api/v1/%s", path)

	if query != nil {
		u.RawQuery = query.Encode()
	}

	return u
}

func trimDomain(domain string) string {
	return strings.TrimSuffix(domain, ".")
}

func makeDomainCanonical(domain string) string {
	return fmt.Sprintf("%s.", trimDomain(domain))
}

func (p *Client) newRequest(ctx context.Context, method string, path string, query *url.Values, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(body)
	}

	apiURL := generateAPIURL(p.Scheme, p.Hostname, p.Port, path, query)
	req, err := http.NewRequestWithContext(ctx, method, apiURL.String(), buf)
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
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    "Unauthorized",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer func() {
			_ = resp.Body.Close()
		}()
		var message string

		if resp.Header.Get("Content-Type") == "application/json" {
			apiError := new(Error)
			_ = json.NewDecoder(resp.Body).Decode(&apiError)
			message = apiError.Message
		} else {
			messageBytes, _ := io.ReadAll(resp.Body)
			message = string(messageBytes)
		}

		return resp, &Error{
			Status:     resp.Status,
			StatusCode: resp.StatusCode,
			Message:    message,
		}
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		defer func() {
			_ = resp.Body.Close()
		}()

		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return resp, err
}
