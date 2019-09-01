package powerdns

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	testBaseURL string = "http://localhost:8080"
	testVhost   string = "localhost"
	testAPIKey  string = "apipw"
)

func generateTestAPIURL() string {
	return fmt.Sprintf("%s/api/v1", testBaseURL)
}

func generateTestAPIVhostURL() string {
	return fmt.Sprintf("%s/servers/%s", generateTestAPIURL(), testVhost)
}

func initialisePowerDNSTestClient() *PowerDNS {
	return NewClient(testBaseURL, testVhost, map[string]string{"X-API-Key": testAPIKey}, nil)
}

func TestError(t *testing.T) {
	myError := &Error{Message: "foo"}
	if myError.Error() != "foo" {
		t.Error("Error method returns invalid format")
	}
}

func TestHandleAPIClientError(t *testing.T) {
	tmpl := &Error{Status: "400 Bad Request", Message: "foo"}
	e := handleAPIClientError(&http.Response{Status: "400 Bad Request", StatusCode: 400}, nil, &Error{Message: "foo"})
	if !reflect.DeepEqual(tmpl, e) {
		t.Error("API client error was not handled properly")
	}
}

func TestNewClient(t *testing.T) {
	tmpl := &PowerDNS{"http", "localhost", "8080", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient}
	p := NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, http.DefaultClient)
	if !reflect.DeepEqual(tmpl, p) {
		t.Error("NewClient returns invalid PowerDNS object")
	}
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
