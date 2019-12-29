package powerdns

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

func generateTestRecord(client *Client, domain string, autoAddRecord bool) string {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("test-%d.%s", rand.Int(), domain)

	if httpmock.Disabled() && autoAddRecord {
		if err := client.Records.Add(domain, name, "TXT", 300, []string{"\"Testing...\""}); err != nil {
			fmt.Printf("Error creating record: %s\n", name)
			fmt.Printf("%s\n", err)
		} else {
			fmt.Printf("Created record %s\n", name)
		}
	}

	return name
}

func registerRecordMockResponder(testDomain string) {
	httpmock.RegisterResponder("PATCH", generateTestAPIVHostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}
			zoneMock := Zone{
				Name: String(fixDomainSuffix(testDomain)),
				URL:  String("/api/v1/servers/" + testVHost + "/zones/" + fixDomainSuffix(testDomain)),
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)
}

func TestAddRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, false)
	registerRecordMockResponder(testDomain)
	if err := p.Records.Add(testDomain, testRecordName, "TXT", 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestAddRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordName, "TXT", 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestChangeRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, true)
	registerRecordMockResponder(testDomain)
	if err := p.Records.Change(testDomain, testRecordName, "TXT", 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Change(testDomain, testRecordName, "TXT", 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteRecord(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, true)
	registerRecordMockResponder(testDomain)
	if err := p.Records.Delete(testDomain, testRecordName, "TXT"); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Delete(testDomain, testRecordName, "TXT"); err == nil {
		t.Error("error is nil")
	}
}
