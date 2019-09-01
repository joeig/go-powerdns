package powerdns

import (
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func generateTestRecord(zone *Zone, autoAddRecord bool) string {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("test-%d.%s", rand.Int(), zone.Name)

	if httpmock.Disabled() && autoAddRecord {
		if err := zone.AddRecord(name, "TXT", 300, []string{"\"Testing...\""}); err != nil {
			fmt.Printf("Error creating record: %s\n", name)
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("Created record %s\n", name)
		}
	}

	return name
}

func TestAddRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	testRecordName := generateTestRecord(z, false)

	httpmock.RegisterResponder("PATCH", generateTestAPIVhostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewBytesResponse(200, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	if z.AddRecord(testRecordName, "TXT", 300, []string{"\"bar\""}) != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	testRecordName := generateTestRecord(z, true)

	httpmock.RegisterResponder("PATCH", generateTestAPIVhostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewBytesResponse(200, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	if z.ChangeRecord(testRecordName, "TXT", 300, []string{"\"bar\""}) != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	testRecordName := generateTestRecord(z, true)

	httpmock.RegisterResponder("PATCH", generateTestAPIVhostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				zoneMock := Zone{
					Name: fixDomainSuffix(testDomain),
					URL:  "/api/v1/servers/" + testVhost + "/zones/" + fixDomainSuffix(testDomain),
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	if z.DeleteRecord(testRecordName, "TXT") != nil {
		t.Errorf("%s", err)
	}
}
