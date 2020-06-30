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

func validateChangeType(changeType ChangeType) error {
	matched, err := regexp.MatchString(`^(REPLACE|DELETE)$`, string(changeType))
	if matched == false || err != nil {
		return &Error{}
	}
	return nil
}

func validateRRType(rrType RRType) error {
	matched, err := regexp.MatchString(`^(A|AAAA|AFSDB|ALIAS|CAA|CERT|CDNSKEY|CDS|CNAME|DNSKEY|DNAME|DS|HINFO|KEY|LOC|MX|NAPTR|NS|NSEC|NSEC3|NSEC3PARAM|OPENPGPKEY|PTR|RP|RRSIG|SOA|SPF|SSHFP|SRV|TKEY|TSIG|TLSA|SMIMEA|TXT|URI)$`, string(rrType))
	if matched == false || err != nil {
		return &Error{}
	}
	return nil
}

func validateCNAMEContent(content string) error {
	if !strings.HasSuffix(content, ".") {
		return &Error{}
	}
	return nil
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

			var rrsets RRsets
			if json.NewDecoder(req.Body).Decode(&rrsets) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			for _, set := range rrsets.Sets {
				if validateChangeType(*set.ChangeType) != nil {
					log.Print("Invalid change type", *set.ChangeType)
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}

				if validateRRType(*set.Type) != nil {
					log.Print("Invalid record type", *set.Type)
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}

				if *set.Type == RRTypeCNAME || *set.Type == RRTypeMX {
					for _, record := range set.Records {
						if validateCNAMEContent(*record.Content) != nil {
							log.Print("CNAME content validation failed")
							return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
						}
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
	registerRecordMockResponder(testDomain)
	testRecordNameTXT := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordNameTXT, RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
	testRecordNameCNAME := generateTestRecord(p, testDomain, false)
	if err := p.Records.Add(testDomain, testRecordNameCNAME, RRTypeCNAME, 300, []string{"foo.tld"}); err != nil {
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

func TestCanonicalResourceRecordValues(t *testing.T) {
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
			canonicalResourceRecordValues(tc.records)

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

func TestFixRRset(t *testing.T) {
	testCases := []struct {
		rrset                     RRset
		wantFixedCanonicalRecords bool
	}{
		{RRset{Type: RRTypePtr(RRTypeMX), Records: []Record{{Content: String("foo.tld")}}}, true},
		{RRset{Type: RRTypePtr(RRTypeCNAME), Records: []Record{{Content: String("foo.tld")}}}, true},
		{RRset{Type: RRTypePtr(RRTypeA), Records: []Record{{Content: String("foo.tld")}}}, false},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			fixRRSet(&tc.rrset)

			if tc.wantFixedCanonicalRecords {
				for j := range tc.rrset.Records {
					isContent := *tc.rrset.Records[j].Content
					wantContent := makeDomainCanonical(*tc.rrset.Records[j].Content)
					if isContent != wantContent {
						t.Errorf("Comparison failed: %s != %s", isContent, wantContent)
					}
				}
			} else {
				for j := range tc.rrset.Records {
					isContent := *tc.rrset.Records[j].Content
					wrongContent := makeDomainCanonical(*tc.rrset.Records[j].Content)
					if isContent == wrongContent {
						t.Errorf("Comparison failed: %s == %s", isContent, wrongContent)
					}
				}
			}
		})
	}
}

func TestPatchRRSets(t *testing.T) {
	testDomain := generateTestZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	testRecordName := generateTestRecord(p, testDomain, true)
	registerRecordMockResponder(testDomain)

	rrSets := RRsets{}
	rrSetName := makeDomainCanonical(testRecordName)
	rrSets.Sets = []RRset{{Name: &rrSetName, Type: RRTypePtr(RRTypeTXT),
		ChangeType: ChangeTypePtr(ChangeTypeDelete)}}

	if err := p.Records.Patch(testDomain, &rrSets); err != nil {
		t.Errorf("%s", err)
	}
}
