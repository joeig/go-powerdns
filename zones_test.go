package powerdns

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestGetZones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zonesMock := []Zone{
					{
						ID:             "example.com.",
						Name:           "example.com.",
						URL:            "/api/v1/servers/localhost/zones/example.com.",
						Kind:           "Native",
						Serial:         1337,
						NotifiedSerial: 1337,
					},
				}
				return httpmock.NewJsonResponse(200, zonesMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := NewClient("http://localhost:8080/", "localhost", "apipw")
	zones, err := p.GetZones()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(zones) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestGetZone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   "example.com.",
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
					Kind: "Native",
					RRsets: []RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster.example.com. 1337 10800 3600 604800 3600",
								},
							},
						},
					},
					Serial:         1337,
					NotifiedSerial: 1337,
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := NewClient("http://localhost:8080/", "localhost", "apipw")
	zone, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != "example.com." {
		t.Error("Received no zone")
	}
}

func TestNotify(t *testing.T) {
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
	httpmock.RegisterResponder("PUT", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/notify",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(200, "{\"result\":\"Notification queued\"}"), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := NewClient("http://localhost:8080/", "localhost", "apipw")
	z, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	notifyResult, err := z.Notify()
	if err != nil {
		t.Errorf("%s", err)
	}
	if notifyResult.Result != "Notification queued" {
		t.Error("Notification was not queued successfully")
	}
}

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

	p := NewClient("http://localhost:8080/", "localhost", "apipw")
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

	p := NewClient("http://localhost:8080/", "localhost", "apipw")
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

	p := NewClient("http://localhost:8080/", "localhost", "apipw")
	z, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if z.DeleteRecord("foo.example.com", "TXT") != nil {
		t.Errorf("%s", err)
	}
}

func TestExport(t *testing.T) {
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
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/export",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(200, "example.com.	3600	SOA	a.misconfigured.powerdns.server. hostmaster.example.com. 1 10800 3600 604800 3600"), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := NewClient("http://localhost:8080/", "localhost", "apipw")
	z, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if _, err := z.Export(); err != nil {
		t.Errorf("%s", err)
	}
}
