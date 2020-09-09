package powerdns

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"testing"
	"time"
)

func generateTestZone(autoAddZone bool) string {
	rand.Seed(time.Now().UnixNano())
	domain := fmt.Sprintf("test-%d.com", rand.Int())

	if httpmock.Disabled() && autoAddZone {
		pdns := initialisePowerDNSTestClient()
		zone, err := pdns.Zones.AddNative(domain, true, "", false, "", "", true, []string{"ns.foo.tld."})
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

func validateZoneType(zoneType ZoneType) error {
	if zoneType != "Zone" {
		return &Error{}
	}
	return nil
}

func validateZoneKind(zoneKind ZoneKind) error {
	matched, err := regexp.MatchString(`^(Native|Master|Slave)$`, string(zoneKind))
	if matched == false || err != nil {
		return &Error{}
	}
	return nil
}

func registerZonesMockResponder() {
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			testDomain := "example.com"
			zonesMock := []Zone{
				{
					ID:             String(makeDomainCanonical(testDomain)),
					Name:           String(makeDomainCanonical(testDomain)),
					URL:            String("/api/v1/servers/" + testVHost + "/zones/" + makeDomainCanonical(testDomain)),
					Kind:           ZoneKindPtr(NativeZoneKind),
					Serial:         Uint32(1337),
					NotifiedSerial: Uint32(1337),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, zonesMock)
		},
	)
}

func registerZoneMockResponder(testDomain string, zoneKind ZoneKind) {
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			zoneMock := Zone{
				ID:   String(makeDomainCanonical(testDomain)),
				Name: String(makeDomainCanonical(testDomain)),
				URL:  String("/api/v1/servers/" + testVHost + "/zones/" + makeDomainCanonical(testDomain)),
				Kind: ZoneKindPtr(NativeZoneKind),
				RRsets: []RRset{
					{
						Name: String(makeDomainCanonical(testDomain)),
						Type: RRTypePtr(RRTypeSOA),
						TTL:  Uint32(3600),
						Records: []Record{
							{
								Content: String("a.misconfigured.powerdns.server. hostmaster." + makeDomainCanonical(testDomain) + " 1337 10800 3600 604800 3600"),
							},
						},
					},
				},
				Serial:         Uint32(1337),
				NotifiedSerial: Uint32(1337),
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)

	httpmock.RegisterResponder("POST", generateTestAPIVHostURL()+"/zones",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var zone Zone
			if json.NewDecoder(req.Body).Decode(&zone) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			if validateZoneType(*zone.Type) != nil {
				log.Print("Invalid zone type", *zone.Type)
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			if validateZoneKind(*zone.Kind) != nil {
				log.Print("Invalid zone kind", *zone.Kind)
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			var zoneMock Zone
			if zoneKind == NativeZoneKind || zoneKind == MasterZoneKind {
				zoneMock = Zone{
					ID:   String(makeDomainCanonical(testDomain)),
					Name: String(makeDomainCanonical(testDomain)),
					Type: ZoneTypePtr(ZoneZoneType),
					URL:  String("api/v1/servers/" + testVHost + "/zones/" + makeDomainCanonical(testDomain)),
					Kind: ZoneKindPtr(zoneKind),
					RRsets: []RRset{
						{
							Name: String(makeDomainCanonical(testDomain)),
							Type: RRTypePtr(RRTypeSOA),
							TTL:  Uint32(3600),
							Records: []Record{
								{
									Content:  String("a.misconfigured.powerdns.server. hostmaster." + makeDomainCanonical(testDomain) + " 0 10800 3600 604800 3600"),
									Disabled: Bool(false),
								},
							},
						},
						{
							Name: String(makeDomainCanonical(testDomain)),
							Type: RRTypePtr(RRTypeNS),
							TTL:  Uint32(3600),
							Records: []Record{
								{
									Content:  String("ns.example.tld."),
									Disabled: Bool(false),
								},
							},
						},
					},
					Serial:      Uint32(0),
					Masters:     []string{},
					DNSsec:      Bool(true),
					Nsec3Param:  String(""),
					Nsec3Narrow: Bool(false),
					SOAEdit:     String("foo"),
					SOAEditAPI:  String("foo"),
					APIRectify:  Bool(true),
					Account:     String(""),
				}
			} else if zoneKind == SlaveZoneKind {
				zoneMock = Zone{
					ID:          String(makeDomainCanonical(testDomain)),
					Name:        String(makeDomainCanonical(testDomain)),
					Type:        ZoneTypePtr(ZoneZoneType),
					URL:         String("api/v1/servers/" + testVHost + "/zones/" + makeDomainCanonical(testDomain)),
					Kind:        ZoneKindPtr(zoneKind),
					Serial:      Uint32(0),
					Masters:     []string{"127.0.0.1"},
					DNSsec:      Bool(true),
					Nsec3Param:  String(""),
					Nsec3Narrow: Bool(false),
					SOAEdit:     String(""),
					SOAEditAPI:  String("DEFAULT"),
					APIRectify:  Bool(true),
					Account:     String(""),
				}
			} else {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Entity"), nil
			}

			return httpmock.NewJsonResponse(http.StatusCreated, zoneMock)
		},
	)

	httpmock.RegisterResponder("PUT", generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewBytesResponse(http.StatusNoContent, []byte{}), nil
		},
	)

	httpmock.RegisterResponder("DELETE", generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewBytesResponse(http.StatusNoContent, []byte{}), nil
		},
	)

	httpmock.RegisterResponder("PUT", generateTestAPIVHostURL()+"/zones/"+testDomain+"/notify",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, "{\"result\":\"Notification queued\"}"), nil
		},
	)

	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+testDomain+"/export",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			acceptHeader := req.Header.Get("Accept")

			if acceptHeader != "" && acceptHeader != "text/html" {
				log.Print("Accept type must be text/html")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewStringResponse(http.StatusOK, makeDomainCanonical(testDomain)+"	3600	SOA	a.misconfigured.powerdns.server. hostmaster."+makeDomainCanonical(testDomain)+" 1 10800 3600 604800 3600"), nil
		},
	)
}

func TestListZones(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZonesMockResponder()

	p := initialisePowerDNSTestClient()
	zones, err := p.Zones.List()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(zones) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestListZonesError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Zones.List(); err == nil {
		t.Error("error is nil")
	}
}

func TestGetZone(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, NativeZoneKind)

	p := initialisePowerDNSTestClient()
	zone, err := p.Zones.Get(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *zone.ID != makeDomainCanonical(testDomain) {
		t.Error("Received no zone")
	}
}

func TestGetZonesError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Zones.Get(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestAddNativeZone(t *testing.T) {
	testDomain := generateTestZone(false)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, NativeZoneKind)

	p := initialisePowerDNSTestClient()
	zone, err := p.Zones.AddNative(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if *zone.ID != makeDomainCanonical(testDomain) || *zone.Kind != NativeZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddNativeZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Zones.AddNative(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."}); err == nil {
		t.Error("error is nil")
	}
}

func TestAddMasterZone(t *testing.T) {
	testDomain := generateTestZone(false)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, MasterZoneKind)

	p := initialisePowerDNSTestClient()
	zone, err := p.Zones.AddMaster(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."})
	if err != nil {
		t.Errorf("%s", err)
	}
	if *zone.ID != makeDomainCanonical(testDomain) || *zone.Kind != MasterZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddMasterZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Zones.AddMaster(testDomain, true, "", false, "foo", "foo", true, []string{"ns.foo.tld."}); err == nil {
		t.Error("error is nil")
	}
}

func TestAddSlaveZone(t *testing.T) {
	testDomain := generateTestZone(false)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, SlaveZoneKind)

	p := initialisePowerDNSTestClient()
	zone, err := p.Zones.AddSlave(testDomain, []string{"127.0.0.1"})
	if err != nil {
		t.Errorf("%s", err)
	}
	if *zone.ID != makeDomainCanonical(testDomain) || *zone.Kind != SlaveZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddSlaveZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Zones.AddSlave(testDomain, []string{"ns5.foo.tld."}); err == nil {
		t.Error("error is nil")
	}
}

func TestAddZone(t *testing.T) {
	testDomain := generateTestZone(false)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, MasterZoneKind)

	p := initialisePowerDNSTestClient()

	z := Zone{
		Name:        String(testDomain),
		Kind:        ZoneKindPtr(MasterZoneKind),
		DNSsec:      Bool(true),
		Nsec3Param:  String(""),
		Nsec3Narrow: Bool(false),
		SOAEdit:     String("foo"),
		SOAEditAPI:  String("foo"),
		APIRectify:  Bool(true),
		Nameservers: []string{"ns.foo.tld."},
	}

	zone, err := p.Zones.Add(&z)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *zone.ID != makeDomainCanonical(testDomain) || *zone.Kind != MasterZoneKind {
		t.Error("Zone wasn't created")
	}
}

func TestAddZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"

	z := Zone{
		Name: String(testDomain),
	}

	if _, err := p.Zones.Add(&z); err == nil {
		t.Error("error is nil")
	}
}

func TestChangeZone(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, NativeZoneKind)

	p := initialisePowerDNSTestClient()

	t.Run("ChangeValidZone", func(t *testing.T) {
		if err := p.Zones.Change(testDomain, &Zone{Nameservers: []string{"ns23.foo.tld."}}); err != nil {
			t.Errorf("%s", err)
		}
	})
	t.Run("ChangeInvalidZone", func(t *testing.T) {
		if err := p.Zones.Change("doesnt-exist", &Zone{Nameservers: []string{"ns23.foo.tld."}}); err == nil {
			t.Errorf("Changing an invalid zone does not return an error")
		}
	})
}

func TestChangeZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if err := p.Zones.Change(testDomain, &Zone{Nameservers: []string{"ns23.foo.tld."}}); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteZone(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, NativeZoneKind)

	p := initialisePowerDNSTestClient()
	if err := p.Zones.Delete(testDomain); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteZoneError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if err := p.Zones.Delete(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestNotify(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, MasterZoneKind)

	p := initialisePowerDNSTestClient()
	notifyResult, err := p.Zones.Notify(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *notifyResult.Result != "Notification queued" {
		t.Error("Notification was not queued successfully")
	}
}

func TestNotifyError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Zones.Notify(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestExport(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain, NativeZoneKind)

	p := initialisePowerDNSTestClient()
	export, err := p.Zones.Export(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if !strings.HasPrefix(string(export), testDomain) {
		t.Errorf("Export payload wrong: \"%s\"", export)
	}
}

func TestExportError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Hostname = "invalid"
	if _, err := p.Zones.Export(testDomain); err == nil {
		t.Error("error is nil")
	}
	p.Port = "x"
	if _, err := p.Zones.Export(testDomain); err == nil {
		t.Error("error is nil")
	}
}
