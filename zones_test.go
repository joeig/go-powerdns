package powerdns

import (
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func generateTestZone(autoAddZone bool) string {
	rand.Seed(time.Now().UnixNano())
	domain := fmt.Sprintf("test-%d.com", rand.Int())

	if httpmock.Disabled() && autoAddZone {
		pdns := NewClient("http://localhost:8080", "localhost", map[string]string{"X-API-Key": "apipw"}, nil)
		zone, err := pdns.AddNativeZone(domain, true, "", false, "", "", true, []string{"ns.foo.tld."})
		if err != nil {
			fmt.Printf("Error creating %s\n", domain)
			fmt.Printf("%v\n", err)
			fmt.Printf("%v\n", zone)
		} else {
			fmt.Printf("Created domain %s\n", domain)
		}
	}

	return domain
}

func registerZoneMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   fixDomainSuffix(testDomain),
					Name: fixDomainSuffix(testDomain),
					URL:  "/api/v1/servers/localhost/zones/" + fixDomainSuffix(testDomain),
					Kind: "Native",
					RRsets: []RRset{
						{
							Name: fixDomainSuffix(testDomain),
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster." + fixDomainSuffix(testDomain) + " 1337 10800 3600 604800 3600",
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
}

func TestGetZones(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zonesMock := []Zone{
					{
						ID:             fixDomainSuffix(testDomain),
						Name:           fixDomainSuffix(testDomain),
						URL:            "/api/v1/servers/localhost/zones/" + fixDomainSuffix(testDomain),
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

	p := initialisePowerDNSTestClient()
	zones, err := p.GetZones()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(zones) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestGetZone(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   fixDomainSuffix(testDomain),
					Name: fixDomainSuffix(testDomain),
					URL:  "/api/v1/servers/localhost/zones/" + fixDomainSuffix(testDomain),
					Kind: NativeZoneKind,
					RRsets: []RRset{
						{
							Name: fixDomainSuffix(testDomain),
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster." + fixDomainSuffix(testDomain) + " 1337 10800 3600 604800 3600",
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

	p := initialisePowerDNSTestClient()
	zone, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != fixDomainSuffix(testDomain) {
		t.Error("Received no zone")
	}
}

func TestAddNativeZone(t *testing.T) {
	testDomain := generateTestZone(false)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:8080/api/v1/servers/localhost/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   fixDomainSuffix(testDomain),
					Name: fixDomainSuffix(testDomain),
					Type: ZoneZoneType,
					URL:  "api/v1/servers/localhost/zones/" + fixDomainSuffix(testDomain),
					Kind: NativeZoneKind,
					RRsets: []RRset{
						{
							Name: fixDomainSuffix(testDomain),
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content:  "a.misconfigured.powerdns.server. hostmaster." + fixDomainSuffix(testDomain) + " 0 10800 3600 604800 3600",
									Disabled: false,
								},
							},
						},
						{
							Name: fixDomainSuffix(testDomain),
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

	p := initialisePowerDNSTestClient()
	zone, err := p.AddNativeZone(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != fixDomainSuffix(testDomain) || zone.Kind != NativeZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddMasterZone(t *testing.T) {
	testDomain := generateTestZone(false)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:8080/api/v1/servers/localhost/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   fixDomainSuffix(testDomain),
					Name: fixDomainSuffix(testDomain),
					Type: ZoneZoneType,
					URL:  "api/v1/servers/localhost/zones/" + fixDomainSuffix(testDomain),
					Kind: MasterZoneKind,
					RRsets: []RRset{
						{
							Name: fixDomainSuffix(testDomain),
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content:  "a.misconfigured.powerdns.server. hostmaster." + fixDomainSuffix(testDomain) + " 0 10800 3600 604800 3600",
									Disabled: false,
								},
							},
						},
						{
							Name: fixDomainSuffix(testDomain),
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

	p := initialisePowerDNSTestClient()
	zone, err := p.AddMasterZone(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != fixDomainSuffix(testDomain) || zone.Kind != MasterZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddSlaveZone(t *testing.T) {
	testDomain := generateTestZone(false)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://localhost:8080/api/v1/servers/localhost/zones",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:          fixDomainSuffix(testDomain),
					Name:        fixDomainSuffix(testDomain),
					Type:        ZoneZoneType,
					URL:         "api/v1/servers/localhost/zones/" + fixDomainSuffix(testDomain),
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

	p := initialisePowerDNSTestClient()
	zone, err := p.AddSlaveZone(testDomain, []string{"ns5.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if zone.ID != fixDomainSuffix(testDomain) || zone.Kind != SlaveZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestChangeZone(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("PUT", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewBytesResponse(204, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()

	t.Run("ChangeValidZone", func(t *testing.T) {
		if err := p.ChangeZone(&Zone{Name: testDomain, Nameservers: []string{"ns23.foo.tld."}}); err != nil {
			t.Errorf("%s", err)
		}
	})
	t.Run("ChangeInvalidZone", func(t *testing.T) {
		if err := p.ChangeZone(&Zone{Name: "", Nameservers: []string{"ns23.foo.tld."}}); err == nil {
			t.Errorf("%s", err)
		}
	})
}

func TestDeleteZone(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("DELETE", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewBytesResponse(204, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	if err := p.DeleteZone(testDomain); err != nil {
		t.Errorf("%s", err)
	}
}

func TestNotify(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)
	httpmock.RegisterResponder("PUT", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain+"/notify",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(200, "{\"result\":\"Notification queued\"}"), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	z, err := p.GetZone(testDomain)
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
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain+"/export",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(200, fixDomainSuffix(testDomain)+"	3600	SOA	a.misconfigured.powerdns.server. hostmaster."+fixDomainSuffix(testDomain)+" 1 10800 3600 604800 3600"), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()
	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if _, err := z.Export(); err != nil {
		t.Errorf("%s", err)
	}
}
