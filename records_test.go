package powerdns

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"testing"
	"time"
)

func generateTestRecord(client *Client, domain string, autoAddRecord bool) string {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("test-%d.%s", rand.Int(), domain)

	if httpmock.Disabled() && autoAddRecord {
		if err := client.Records.Add(domain, name, RRTypeTXT, 300, []string{"\"Testing...\""}); err != nil {
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

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var requestBodyRRsets RRsets
			err := json.NewDecoder(req.Body).Decode(&requestBodyRRsets)
			if err != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			for _, set := range requestBodyRRsets.Sets {
				if *set.Type != RRTypeCNAME {
					break
				}
				for _, record := range set.Records {
					if !strings.HasSuffix(*record.Content, ".") {
						log.Print("CNAME content validation failed")
						return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
					}
				}
			}

			zoneMock := Zone{
				Name: String(makeDomainCanonical(testDomain)),
				URL:  String("/api/v1/servers/" + testVHost + "/zones/" + makeDomainCanonical(testDomain)),
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
	if err := p.Records.Add(testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
	if err := p.Records.Add(testDomain, testRecordName, RRTypeCNAME, 300, []string{"foo.tld"}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestAddRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
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
	if err := p.Records.Change(testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Change(testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
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
	if err := p.Records.Delete(testDomain, testRecordName, RRTypeTXT); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateTestZone(false)
	testRecordName := generateTestRecord(p, testDomain, false)
	if err := p.Records.Delete(testDomain, testRecordName, RRTypeTXT); err == nil {
		t.Error("error is nil")
	}
}

func TestFixCNAMEResourceRecordValues(t *testing.T) {
	testCases := []struct {
		records     []Record
		wantContent []string
	}{
		{[]Record{{Content: String("foo.tld")}}, []string{"foo.tld."}},
		{[]Record{{Content: String("foo.tld.")}}, []string{"foo.tld."}},
		{[]Record{{Content: String("foo.tld")}, {Content: String("foo.tld.")}}, []string{"foo.tld.", "foo.tld."}},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			fixCNAMEResourceRecordValues(tc.records)
			for j := range tc.records {
				isContent := *tc.records[j].Content
				wantContent := tc.wantContent[j]
				if isContent != wantContent {
					t.Errorf("Comparison failed: %s != %s", isContent, wantContent)
				}
			}
		})
	}
}
