// Package powerdns is a Go client library for the PowerDNS API.
// It's a community project and not associated with the official PowerDNS product itself.
package powerdns

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
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

// WithAPIKey is an option for New to set the API key.
func WithAPIKey(key string) NewOption {
	return func(client *Client) {
		client.apiKey = &key
	}
}

type service struct {
	client *Client
}

// Client configuration structure
type Client struct {
	BaseURL string
	VHost   string
	Headers map[string]string

	httpClient *http.Client
	apiKey     *string

	common service // Reuse a single struct instead of allocating one for each service on the heap

	Config     *ConfigService
	Cryptokeys *CryptokeysService
	Metadata   *MetadataService
	Records    *RecordsService
	Servers    *ServersService
	Statistics *StatisticsService
	Zones      *ZonesService
	// Deprecated: Use TSIGKeys instead. TSIGKey will be removed with the next major version.
	TSIGKey  *TSIGKeysService
	TSIGKeys *TSIGKeysService
}

// New initializes a new client instance.
func New(baseURL string, vHost string, options ...NewOption) *Client {
	client := &Client{
		BaseURL:    baseURL,
		VHost:      parseVHost(vHost),
		httpClient: http.DefaultClient,
	}

	client.common.client = client

	client.Config = (*ConfigService)(&client.common)
	client.Cryptokeys = (*CryptokeysService)(&client.common)
	client.Metadata = (*MetadataService)(&client.common)
	client.Records = (*RecordsService)(&client.common)
	client.Servers = (*ServersService)(&client.common)
	client.Statistics = (*StatisticsService)(&client.common)
	client.Zones = (*ZonesService)(&client.common)
	client.TSIGKeys = (*TSIGKeysService)(&client.common)
	client.TSIGKey = client.TSIGKeys

	for _, option := range options {
		option(client)
	}

	return client
}

func parseVHost(vHost string) string {
	if vHost == "" {
		return "localhost"
	}
	return vHost
}

func generateAPIURL(baseURL, pathFragment string, query *url.Values) (*url.URL, error) {
	newURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	newURL.Path = path.Join("/api/v1/", pathFragment)

	if query != nil {
		newURL.RawQuery = query.Encode()
	}

	return newURL, nil
}

func trimDomain(domain string) string {
	return strings.TrimSuffix(domain, ".")
}

func makeDomainCanonical(domain string) string {
	return fmt.Sprintf("%s.", trimDomain(domain))
}

func (p *Client) newRequest(ctx context.Context, method string, pathFragment string, query *url.Values, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(body)
	}

	apiURL, err := generateAPIURL(p.BaseURL, pathFragment, query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, apiURL.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "go-powerdns")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")
	}

	if p.apiKey != nil {
		req.Header.Set("X-API-Key", *p.apiKey)
	}

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

		if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
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
