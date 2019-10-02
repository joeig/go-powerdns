package powerdns

import (
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("TestValidURL", func(t *testing.T) {
		tmpl := &Client{"http", "localhost", "8080", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient, service{}, nil, nil, nil, nil, nil}
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
		errors := []string{}
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
	if _, err := p.newRequest("GET", "servers", nil); err != nil {
		t.Error("NewRequest returns an error")
	}
}

func TestDo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/return-401", generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/return-404", generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, "Not Found"), nil
		},
	)
	httpmock.RegisterResponder("GET", fmt.Sprintf("%s/server", generateTestAPIURL()),
		func(req *http.Request) (*http.Response, error) {
			mock := Error{
				Status:  "Not Found",
				Message: "Not Found",
			}
			return httpmock.NewJsonResponse(http.StatusNotImplemented, mock)
		},
	)
	p := initialisePowerDNSTestClient()

	t.Run("Test401Handling", func(t *testing.T) {
		req, _ := p.newRequest("GET", "return-401", nil)
		if _, err := p.do(req, nil); err == nil {
			t.Error("401 response does not result into an error")
		}
	})
	t.Run("Test404Handling", func(t *testing.T) {
		req, _ := p.newRequest("GET", "return-404", nil)
		if _, err := p.do(req, nil); err == nil {
			t.Error("404 response does not result into an error")
		}
	})
	t.Run("TestJSONResponseHandling", func(t *testing.T) {
		req, _ := p.newRequest("GET", "server", &Server{})
		if _, err := p.do(req, nil); err.(*Error).Message != "Not Found" {
			t.Error("501 JSON response does not result into Error structure")
		}
	})
}

func TestParseBaseURL(t *testing.T) {
	testCases := []struct {
		baseURL      string
		wantScheme   string
		wantHostname string
		wantPort     string
	}{
		{"https://example.com", "https", "example.com", "443"},
		{"http://example.com", "http", "example.com", "80"},
		{"https://example.com:8080", "https", "example.com", "8080"},
		{"http://example.com:8080", "http", "example.com", "8080"},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			scheme, hostname, port, err := parseBaseURL(tc.baseURL)
			if err != nil {
				t.Errorf("%s is not a valid url: %v", tc.baseURL, err)
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

	t.Run("InvalidURL", func(t *testing.T) {
		if _, _, _, err := parseBaseURL("http%%%foo"); err == nil {
			t.Error("Invalid URL does not return an error")
		}
	})
}

func TestParseVhost(t *testing.T) {
	t.Run("ValidVhost", func(t *testing.T) {
		if parseVhost("example.com") != "example.com" {
			t.Error("Valid vhost returned invalid value")
		}
	})
	t.Run("MissingVhost", func(t *testing.T) {
		if parseVhost("") != "localhost" {
			t.Error("Missing vhost did not return localhost")
		}
	})
}

func TestGenerateAPIURL(t *testing.T) {
	tmpl := "https://localhost:8080/api/v1/foo"
	g := generateAPIURL("https", "localhost", "8080", "foo")
	if tmpl != g.String() {
		t.Errorf("Template does not match generated API URL: %s", g.String())
	}
}
