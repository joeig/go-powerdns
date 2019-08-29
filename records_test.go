package powerdns

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestAddRecord(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("PATCH", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if z.AddRecord("foo.example.com", "TXT", 300, []string{"bar"}) != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecord(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("PATCH", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if z.ChangeRecord("foo.example.com", "TXT", 300, []string{"bar"}) != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("PATCH", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if z.DeleteRecord("foo.example.com", "TXT") != nil {
		t.Errorf("%s", err)
	}
}
