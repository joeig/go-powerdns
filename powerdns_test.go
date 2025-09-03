package powerdns

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"testing"

	"github.com/jarcoal/httpmock"
)

const (
	testBaseURL string = "http://localhost:8080"
	testVHost   string = "localhost"
	testAPIKey  string = "apipw"
)

func generateTestAPIURL() string {
	return fmt.Sprintf("%s/api/v1", testBaseURL)
}

func generateTestAPIVHostURL() string {
	return fmt.Sprintf("%s/servers/%s", generateTestAPIURL(), testVHost)
}

func verifyAPIKey(req *http.Request) *http.Response {
	if req.Header.Get("X-Api-Key") != testAPIKey {
		return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized")
	}
	return nil
}
func initialisePowerDNSTestClient() *Client {
	client := New(testBaseURL, testVHost, WithAPIKey(testAPIKey))
	return client
}

func registerDoMockResponder() {
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/servers/doesnt-exist", generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}
			return httpmock.NewStringResponse(http.StatusNotFound, "Not Found"), nil
		},
	)

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/servers/localhost", generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			return verifyAPIKey(req), nil
		},
	)

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/server", generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			mock := Error{
				Status:     "Not Found",
				StatusCode: http.StatusNotFound,
				Message:    "Not Found",
			}
			return httpmock.NewJsonResponse(http.StatusNotImplemented, mock)
		},
	)
}

func TestWithHeaders(t *testing.T) {
	p := &Client{}
	withHeaders := WithHeaders(map[string]string{"X-Test-Header": "test-header"})
	withHeaders(p)
	if !maps.Equal(p.Headers, map[string]string{"X-Test-Header": "test-header"}) {
		t.Error("Unexpected header")
	}
}

func TestWithHttpClient(t *testing.T) {
	p := &Client{}
	httpClient := &http.Client{}
	withHTTPClient := WithHTTPClient(httpClient)
	withHTTPClient(p)
	if p.httpClient != httpClient {
		t.Error("Unexpected HTTP client")
	}
}

func TestWithAPIKey(t *testing.T) {
	p := &Client{}
	withAPIKey := WithAPIKey("apipw")
	withAPIKey(p)
	if *p.apiKey != "apipw" {
		t.Error("Unexpected API key")
	}
}

func TestNew(t *testing.T) {
	t.Run("TestNoOptions", func(t *testing.T) {
		p := New("http://localhost:8080", "localhost")
		if p.BaseURL != "http://localhost:8080" {
			t.Error("New returns invalid base URL")
		}
		if p.VHost != "localhost" {
			t.Error("New returns invalid vHost")
		}
		if !maps.Equal(p.Headers, map[string]string{}) {
			t.Error("New returns invalid headers")
		}
		if p.httpClient != http.DefaultClient {
			t.Error("New returns invalid HTTP client")
		}
		if p.common.client != p {
			t.Error("New returns invalid common client")
		}
	})

	t.Run("TestOptionInvocation", func(t *testing.T) {
		testOptionInvocationCount := 0
		testOption := func(client *Client) {
			testOptionInvocationCount++
		}
		_ = New("http://localhost:8080", "localhost", testOption, testOption)

		if testOptionInvocationCount != 2 {
			t.Error("New does not call all options")
		}
	})
}

func TestNewRequest(t *testing.T) {
	t.Run("TestValidRequest", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		if _, err := p.newRequest(context.Background(), http.MethodGet, "servers", nil, nil); err != nil {
			t.Error("error is not nil")
		}
	})

	t.Run("TestUserAgentHeader", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers", nil, nil)
		if req.Header.Get("User-Agent") != "go-powerdns" {
			t.Error("Unexpected user agent header")
		}
	})

	t.Run("TestContentTypeHeaderWithoutBody", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers", nil, nil)
		if req.Header.Get("Content-Type") != "" {
			t.Error("Unexpected content type header")
		}
		if req.Header.Get("Accept") != "" {
			t.Error("Unexpected accept header")
		}
	})

	t.Run("TestContentTypeHeaderWithBody", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers", nil, "test-body")
		if req.Header.Get("Content-Type") != "application/json" {
			t.Error("Unexpected content type header")
		}
		if req.Header.Get("Accept") != "application/json" {
			t.Error("Unexpected accept header")
		}
	})

	t.Run("TestAPIKeyHeader", func(t *testing.T) {
		p := New(testBaseURL, testVHost, WithAPIKey("test-key"))
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers", nil, nil)
		if req.Header.Get("X-API-Key") != "test-key" {
			t.Error("Unexpected API key header")
		}
	})

	t.Run("TestCustomHeaders", func(t *testing.T) {
		p := New(testBaseURL, testVHost, WithHeaders(map[string]string{"X-Test-Header": "test-header"}))
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers", nil, nil)
		if req.Header.Get("X-Test-Header") != "test-header" {
			t.Error("Unexpected API key header")
		}
	})

	t.Run("TestInvalidMethod", func(t *testing.T) {
		p := New(testBaseURL, testVHost, WithHeaders(map[string]string{"X-Test-Header": "test-header"}))
		_, err := p.newRequest(context.Background(), " ", "servers", nil, nil)
		if err == nil {
			t.Error("err expected")
		}
	})
}

func TestDo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerDoMockResponder()

	t.Run("TestStringErrorResponse", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers/doesnt-exist", nil, nil)
		if _, err := p.do(req, nil); err == nil {
			t.Error("err is nil")
		}
	})
	t.Run("Test401Handling", func(t *testing.T) {
		p := New(testBaseURL, testVHost)
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers/localhost", nil, nil)
		if _, err := p.do(req, nil); err.Error() != "Unauthorized" {
			t.Error("401 response does not result into an error with correct message.")
		}
	})
	t.Run("TestErrorHandling", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest(context.Background(), http.MethodGet, "servers/doesnt-exist", nil, nil)
		_, err := p.do(req, nil)
		wantResultBeforePowerDNSAuth49 := "Not Found"
		wantResultFromPowerDNSAuth49 := "Method Not Allowed"
		if err.Error() != wantResultBeforePowerDNSAuth49 && err.Error() != wantResultFromPowerDNSAuth49 {
			t.Error("Error response does not result into an error with correct message.", err.Error())
		}
	})
	t.Run("TestJSONErrorHandling", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest(context.Background(), http.MethodGet, "server", nil, nil)
		_, err := p.do(req, nil)
		wantResultBeforePowerDNSAuth49 := "Not Found"
		wantResultFromPowerDNSAuth49 := "Method Not Allowed"
		if err.Error() != wantResultBeforePowerDNSAuth49 && err.Error() != wantResultFromPowerDNSAuth49 {
			t.Error("Error response does not result into an error with correct message.", err.Error())
		}
	})
}

func TestParseVHost(t *testing.T) {
	testCases := []struct {
		vHost     string
		wantVHost string
	}{
		{"example.com", "example.com"},
		{"", "localhost"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			if parseVHost(tc.vHost) != tc.wantVHost {
				t.Error("parseVHost returned an invalid value")
			}
		})
	}
}

func TestGenerateAPIURL(t *testing.T) {
	testQuery := &url.Values{}
	testQuery.Add("a", "b")
	testCases := []struct {
		baseURL string
		path    string
		query   *url.Values
		wantURL string
		wantErr *error
	}{
		{"https://localhost:8080", "foo", testQuery, "https://localhost:8080/api/v1/foo?a=b", nil},
		{"http://localhost:8080", "foo", testQuery, "http://localhost:8080/api/v1/foo?a=b", nil},
		{"http://localhost:1337", "foo", testQuery, "http://localhost:1337/api/v1/foo?a=b", nil},
		{"https://127.1.2.3:8080", "foo", testQuery, "https://127.1.2.3:8080/api/v1/foo?a=b", nil},
		{"https://[fd06:4c9a:99b0::1]:8080", "foo", testQuery, "https://[fd06:4c9a:99b0::1]:8080/api/v1/foo?a=b", nil},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			newURL, err := generateAPIURL(tc.baseURL, tc.path, tc.query)
			if tc.wantURL != newURL.String() {
				t.Errorf("generateAPIURL returned an invalid value: %q", newURL.String())
			}
			if tc.wantErr != nil && errors.Is(err, *tc.wantErr) {
				t.Errorf("generateAPIURL returned an invalid error: %q", err)
			}
		})
	}
}

func TestTrimDomain(t *testing.T) {
	testCases := []struct {
		domain     string
		wantDomain string
	}{
		{"example.com.", "example.com"},
		{"example.com", "example.com"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			if trimDomain(tc.domain) != tc.wantDomain {
				t.Error("trimDomain returned an invalid value")
			}
		})
	}
}

func TestMakeDomainCanonical(t *testing.T) {
	testCases := []struct {
		domain     string
		wantDomain string
	}{
		{"example.com.", "example.com."},
		{"example.com", "example.com."},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			if makeDomainCanonical(tc.domain) != tc.wantDomain {
				t.Error("makeDomainCanonical returned an invalid value")
			}
		})
	}
}
