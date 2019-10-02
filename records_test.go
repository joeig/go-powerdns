package powerdns

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestAddRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, false)

	httpmock.RegisterResponder("PATCH", generateTestAPIVhostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewBytesResponse(200, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	if err := p.Records.Add(testDomain, testRecordName, "TXT", 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, true)

	httpmock.RegisterResponder("PATCH", generateTestAPIVhostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewBytesResponse(200, []byte{}), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	if err := p.Records.Change(testDomain, testRecordName, "TXT", 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, true)

	httpmock.RegisterResponder("PATCH", generateTestAPIVhostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				zoneMock := Zone{
					Name: String(fixDomainSuffix(testDomain)),
					URL:  String("/api/v1/servers/" + testVhost + "/zones/" + fixDomainSuffix(testDomain)),
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	if err := p.Records.Delete(testDomain, testRecordName, "TXT"); err != nil {
		t.Errorf("%s", err)
	}
}
