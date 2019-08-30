package powerdns

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
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
						Kind:           NativeZoneKind,
						Serial:         1337,
						NotifiedSerial: 1337,
					},
				}
				return httpmock.NewJsonResponse(200, zonesMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
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
					Kind: NativeZoneKind,
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

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	zone, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != "example.com." {
		t.Error("Received no zone")
	}
}

func TestAddNativeZone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:8080/api/v1/servers/localhost/zones/",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   "example.com.",
					Name: "example.com.",
					Type: ZoneZoneType,
					URL:  "api/v1/servers/localhost/zones/example.com.",
					Kind: NativeZoneKind,
					RRsets: []RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content:  "a.misconfigured.powerdns.server. hostmaster.example.com. 0 10800 3600 604800 3600",
									Disabled: false,
								},
							},
						},
						{
							Name: "example.com.",
							Type: "NS",
							TTL:  3600,
							Records: []Record{
								{
									Content:  "ns.example.tld.",
									Disabled: false,
								},
							},
						},
					},
					Serial:      0,
					Masters:     []string{},
					DNSsec:      true,
					Nsec3Param:  "",
					Nsec3Narrow: false,
					SOAEdit:     "foo",
					SOAEditAPI:  "foo",
					APIRectify:  true,
					Account:     "",
				}
				return httpmock.NewJsonResponse(201, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	zone, err := p.AddNativeZone("example.com", true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != "example.com." || zone.Kind != NativeZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddMasterZone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:8080/api/v1/servers/localhost/zones/",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   "example.com.",
					Name: "example.com.",
					Type: ZoneZoneType,
					URL:  "api/v1/servers/localhost/zones/example.com.",
					Kind: MasterZoneKind,
					RRsets: []RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content:  "a.misconfigured.powerdns.server. hostmaster.example.com. 0 10800 3600 604800 3600",
									Disabled: false,
								},
							},
						},
						{
							Name: "example.com.",
							Type: "NS",
							TTL:  3600,
							Records: []Record{
								{
									Content:  "ns.example.tld.",
									Disabled: false,
								},
							},
						},
					},
					Serial:      0,
					Masters:     []string{},
					DNSsec:      true,
					Nsec3Param:  "",
					Nsec3Narrow: false,
					SOAEdit:     "foo",
					SOAEditAPI:  "foo",
					APIRectify:  true,
					Account:     "",
				}
				return httpmock.NewJsonResponse(201, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	zone, err := p.AddMasterZone("example.com", true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != "example.com." || zone.Kind != MasterZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddSlaveZone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:8080/api/v1/servers/localhost/zones/",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:          "example.com.",
					Name:        "example.com.",
					Type:        ZoneZoneType,
					URL:         "api/v1/servers/localhost/zones/example.com.",
					Kind:        SlaveZoneKind,
					Serial:      0,
					Masters:     []string{"ns5.foo.tld."},
					DNSsec:      true,
					Nsec3Param:  "",
					Nsec3Narrow: false,
					SOAEdit:     "",
					SOAEditAPI:  "DEFAULT",
					APIRectify:  true,
					Account:     "",
				}
				return httpmock.NewJsonResponse(201, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	zone, err := p.AddSlaveZone("example.com", []string{"ns5.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != "example.com." || zone.Kind != SlaveZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestChangeZone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("PUT", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewBytesResponse(204, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	if err := p.ChangeZone(&Zone{Name: "example.com", Nameservers: []string{"ns23.foo.tld."}}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteZone(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("DELETE", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewBytesResponse(204, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	if err := p.DeleteZone("example.com"); err != nil {
		t.Errorf("%s", err)
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

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
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

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	if err != nil {
		t.Errorf("%s", err)
	}
	if _, err := z.Export(); err != nil {
		t.Errorf("%s", err)
	}
}
