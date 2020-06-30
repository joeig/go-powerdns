package powerdns

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"net/url"
	"testing"
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
	return NewClient(testBaseURL, testVHost, map[string]string{"X-API-Key": testAPIKey}, nil)
}

func registerDoMockResponder() {

	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/servers/doesntExist", generateTestAPIURL()),
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

func TestNewClient(t *testing.T) {
	t.Run("TestValidURL", func(t *testing.T) {
		tmpl := &Client{"http", "localhost", "8080", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient, service{}, nil, nil, nil, nil, nil, nil}
		p := NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient)
		if p.Hostname != tmpl.Hostname {
			t.Error("NewClient returns invalid Client object")
		}
	})

	t.Run("TestInvalidURL", func(t *testing.T) {
		originalLogFatalf := logFatalf
		defer func() {
			logFatalf = originalLogFatalf
		}()
		var errors []string
		logFatalf = func(format string, args ...interface{}) {
			if len(args) > 0 {
				errors = append(errors, fmt.Sprintf(format, args))
			} else {
				errors = append(errors, format)
			}
		}

		_ = NewClient("http://1.2:foo", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient)

		if len(errors) < 1 {
			t.Error("NewClient does not exit with fatal error")
		}
	})
}

func TestNewRequest(t *testing.T) {
	p := initialisePowerDNSTestClient()

	t.Run("TestValidRequest", func(t *testing.T) {
		if _, err := p.newRequest("GET", "servers", nil, nil); err != nil {
			t.Error("error is not nil")
		}
	})
}

func TestDo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerDoMockResponder()

	t.Run("TestStringErrorResponse", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest("GET", "servers/doesntExist", nil, nil)
		if _, err := p.do(req, nil); err == nil {
			t.Error("err is nil")
		}
	})
	t.Run("Test401Handling", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		p.Headers = nil
		req, _ := p.newRequest("GET", "servers/localhost", nil, nil)
		if _, err := p.do(req, nil); err.Error() != "Unauthorized" {
			t.Error("401 response does not result into an error with correct message.")
		}
	})
	t.Run("Test404Handling", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest("GET", "servers/doesntExist", nil, nil)
		if _, err := p.do(req, nil); err.Error() != "Not Found" {
			t.Error("404 response does not result into an error with correct message.")
		}
	})
	t.Run("TestJSONResponseHandling", func(t *testing.T) {
		p := initialisePowerDNSTestClient()
		req, _ := p.newRequest("GET", "server", nil, &Server{})
		if _, err := p.do(req, nil); err.Error() != "Not Found" {
			t.Error("501 JSON response does not result into an error with correct message.")
		}
	})
}

func TestParseBaseURL(t *testing.T) {
	testCases := []struct {
		baseURL      string
		wantScheme   string
		wantHostname string
		wantPort     string
		wantError    bool
	}{
		{"https://example.com", "https", "example.com", "443", false},
		{"http://example.com", "http", "example.com", "80", false},
		{"https://example.com:8080", "https", "example.com", "8080", false},
		{"http://example.com:8080", "http", "example.com", "8080", false},
		{"http%%%foo", "http", "", "", true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			scheme, hostname, port, err := parseBaseURL(tc.baseURL)

			if err != nil && tc.wantError == true {
				return
			}
			if err != nil && tc.wantError == false {
				t.Error("Error was returned unexpectedly")
			}
			if err == nil && tc.wantError == true {
				t.Error("No error was returned")
			}
			if scheme != tc.wantScheme {
				t.Errorf("Scheme parsing failed: %s != %s", scheme, tc.wantScheme)
			}
			if hostname != tc.wantHostname {
				t.Errorf("Hostname parsing failed: %s != %s", hostname, tc.wantHostname)
			}
			if port != tc.wantPort {
				t.Errorf("Port parsing failed: %s != %s", port, tc.wantPort)
			}
		})
	}
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
	tmpl := "https://localhost:8080/api/v1/foo?a=b"
	query := url.Values{}
	query.Add("a", "b")
	g := generateAPIURL("https", "localhost", "8080", "foo", &query)
	if tmpl != g.String() {
		t.Errorf("Template does not match generated API URL: %s", g.String())
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
