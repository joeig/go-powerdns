package powerdns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

const (
	testARecord    = "10.0.0.1"
	testAAAARecord = "d96e:a60a:99c2:a3db:7b08:c36:3dc1:6d4a"
	testTXTRecord  = "\"Testing...\""
)

var (
	testRecordTXT = record{
		Type:    RRTypeTXT,
		TTL:     uint32(300),
		Content: []string{testTXTRecord},
	}

	testRecordA = record{
		Type:    RRTypeA,
		TTL:     uint32(300),
		Content: []string{testARecord},
	}

	testRecordAAAA = record{
		Type:    RRTypeAAAA,
		TTL:     uint32(300),
		Content: []string{testAAAARecord},
	}
)

type record struct {
	Type    RRType
	TTL     uint32
	Content []string
}

func generateTestRecord(client *Client, domain string, autoAddRecord bool, records ...record) string {
	name := fmt.Sprintf("test-%d.%s", rand.Int(), domain)

	if httpmock.Disabled() && autoAddRecord {
		for _, rec := range records {
			if err := client.Records.Add(context.Background(), domain, name, rec.Type, rec.TTL, rec.Content); err != nil {
				log.Printf("Error creating record: %s: %v\n", name, err)
				continue
			}

			log.Printf("Created record %s\n", name)
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

func registerRecordMockResponder(testDomain, testRecord string) {
	testDomainCanonical := makeDomainCanonical(testDomain)
	testRecordCanonical := makeDomainCanonical(testRecord)

	httpmock.RegisterResponder("PATCH", generateTestAPIVHostURL()+"/zones/"+testDomainCanonical,
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
				Name: String(testDomainCanonical),
				URL:  String("/api/v1/servers/" + testVHost + "/zones/" + testDomainCanonical),
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)
	httpmock.RegisterResponder(http.MethodGet, generateTestAPIVHostURL()+"/zones/"+testDomainCanonical+"?rrset_name="+testRecordCanonical+"&rrset_type=TXT",
		func(req *http.Request) (*http.Response, error) {

			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			zoneMock := Zone{
				Name: String(testDomainCanonical),
				URL:  String("/api/v1/servers/" + testVHost + "/zones/" + testDomainCanonical),
				RRsets: []RRset{
					{
						Name: String(testRecordCanonical),
						Type: RRTypePtr(RRTypeTXT),
						TTL:  Uint32(300),
						Records: []Record{
							{
								Content: String(testTXTRecord),
							},
						},
					},
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)
	httpmock.RegisterResponder(http.MethodGet, generateTestAPIVHostURL()+"/zones/"+testDomainCanonical+"?rrset_name="+testRecordCanonical,
		func(req *http.Request) (*http.Response, error) {

			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			zoneMock := Zone{
				Name: String(testDomainCanonical),
				URL:  String("/api/v1/servers/" + testVHost + "/zones/" + testDomainCanonical),
				RRsets: []RRset{
					{
						Name: String(testRecordCanonical),
						Type: RRTypePtr(RRTypeTXT),
						TTL:  Uint32(300),
						Records: []Record{
							{
								Content: String(testTXTRecord),
							},
						},
					},
					{
						Name: String(testRecordCanonical),
						Type: RRTypePtr(RRTypeA),
						TTL:  Uint32(300),
						Records: []Record{
							{
								Content: String(testARecord),
							},
						},
					},
					{
						Name: String(testRecordCanonical),
						Type: RRTypePtr(RRTypeAAAA),
						TTL:  Uint32(300),
						Records: []Record{
							{
								Content: String(testAAAARecord),
							},
						},
					},
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)
	httpmock.RegisterResponder(http.MethodGet, generateTestAPIVHostURL()+"/zones/"+testDomainCanonical+"?rrset_name="+makeDomainCanonical(testRecord+"notfound"),
		func(req *http.Request) (*http.Response, error) {

			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}
			zoneMock := Zone{
				Name:   String(testDomainCanonical),
				URL:    String("/api/v1/servers/" + testVHost + "/zones/" + testDomainCanonical),
				RRsets: []RRset{},
			}
			return httpmock.NewJsonResponse(http.StatusOK, zoneMock)
		},
	)
}

func TestAddRecord(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	registerRecordMockResponder(testDomain, "")
	testRecordNameTXT := generateTestRecord(p, testDomain, false, testRecordTXT)
	if err := p.Records.Add(context.Background(), testDomain, testRecordNameTXT, RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
	testRecordNameCNAME := generateTestRecord(p, testDomain, false, testRecordTXT)
	if err := p.Records.Add(context.Background(), testDomain, testRecordNameCNAME, RRTypeCNAME, 300, []string{"foo.tld"}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestAddRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateNativeZone(false)
	testRecordName := generateTestRecord(p, testDomain, false, testRecordTXT)
	if err := p.Records.Add(context.Background(), testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestChangeRecord(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, true, testRecordTXT)
	registerRecordMockResponder(testDomain, testRecordName)
	if err := p.Records.Change(context.Background(), testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecordComment(t *testing.T) {
	comment := Comment{
		Content:    String("Example comment"),
		Account:    String("example account"),
		ModifiedAt: Uint64(uint64(time.Now().Unix())),
	}
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, true, testRecordTXT)
	registerRecordMockResponder(testDomain, testRecordName)
	if err := p.Records.Change(context.Background(), testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}, WithComments(comment)); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateNativeZone(false)
	testRecordName := generateTestRecord(p, testDomain, false, testRecordTXT)
	if err := p.Records.Change(context.Background(), testDomain, testRecordName, RRTypeTXT, 300, []string{"\"bar\""}); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteRecord(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	testRecordName := generateTestRecord(p, testDomain, true, testRecordTXT)
	registerRecordMockResponder(testDomain, testRecordName)
	if err := p.Records.Delete(context.Background(), testDomain, testRecordName, RRTypeTXT); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateNativeZone(false)
	testRecordName := generateTestRecord(p, testDomain, false, testRecordTXT)
	if err := p.Records.Delete(context.Background(), testDomain, testRecordName, RRTypeTXT); err == nil {
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

func TestGetRecord(t *testing.T) {
	testDomain := generateNativeZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	testRecordName := generateTestRecord(p, testDomain, true, testRecordTXT, testRecordA, testRecordAAAA)
	registerRecordMockResponder(testDomain, testRecordName)

	testRecordNameCanonical := makeDomainCanonical(testRecordName)

	testCases := []struct {
		testDesc       string
		testRecordName string
		testRecordType *RRType
		expectRRset    []RRset
	}{
		{
			testDesc:       "Get with rrset_name and rrset_type",
			testRecordName: testRecordNameCanonical,
			testRecordType: RRTypePtr(RRTypeTXT),
			expectRRset: []RRset{
				{
					Name: String(testRecordNameCanonical),
					Type: RRTypePtr(RRTypeTXT),
					TTL:  Uint32(300),
					Records: []Record{
						{
							Content: String(testTXTRecord),
						},
					},
				},
			},
		},
		{
			testDesc:       "Get with rrset_name",
			testRecordName: testRecordNameCanonical,
			testRecordType: nil,
			expectRRset: []RRset{
				{
					Name: String(testRecordNameCanonical),
					Type: RRTypePtr(RRTypeA),
					TTL:  Uint32(300),
					Records: []Record{
						{
							Content: String(testARecord),
						},
					},
				},
				{
					Name: String(testRecordNameCanonical),
					Type: RRTypePtr(RRTypeAAAA),
					TTL:  Uint32(300),
					Records: []Record{
						{
							Content: String(testAAAARecord),
						},
					},
				},
				{
					Name: String(testRecordNameCanonical),
					Type: RRTypePtr(RRTypeTXT),
					TTL:  Uint32(300),
					Records: []Record{
						{
							Content: String(testTXTRecord),
						},
					},
				},
			},
		},
		{
			testDesc:       "Get with rrset_name not found",
			testRecordName: makeDomainCanonical(testRecordName + "notfound"),
			testRecordType: nil,
			expectRRset:    nil,
		},
	}

	for n, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d - %s", n, tc.testDesc), func(t *testing.T) {

			fmt.Println("Get ", tc.testRecordName)
			rrsets, err := p.Records.Get(context.Background(), testDomain, tc.testRecordName, tc.testRecordType)
			if err != nil {
				t.Errorf("Unexpected error got %v", err)
			}

			// Sort for consistency
			sort.Slice(rrsets, func(i, j int) bool {
				return *rrsets[i].Type < *rrsets[j].Type
			})

			for i, r := range rrsets {
				if *r.Name != *tc.expectRRset[i].Name {
					t.Errorf("Comparison rrset name failed %v != %v", *r.Name, *tc.expectRRset[i].Name)
				}

				if *r.Type != *tc.expectRRset[i].Type {
					t.Errorf("Comparison rrset type failed %v != %v", *r.Type, *tc.expectRRset[i].Type)
				}

				if *r.TTL != *tc.expectRRset[i].TTL {
					t.Errorf("Comparison rrset TTL failed %v != %v", *r.TTL, *tc.expectRRset[i].TTL)
				}

				for j, rec := range r.Records {
					if *rec.Content != *tc.expectRRset[i].Records[j].Content {
						t.Errorf("Comparison rrset record content failed %v != %v", *rec.Content, *tc.expectRRset[i].Records[j].Content)
					}
				}
			}
		})
	}
}

func TestGetRecordError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	testDomain := generateNativeZone(false)
	testRecordName := generateTestRecord(p, testDomain, false, testRecordTXT)
	if _, err := p.Records.Get(context.Background(), testDomain, testRecordName, RRTypePtr(RRTypeTXT)); err == nil {
		t.Error("error is nil")
	}
}

func TestPatchRRSets(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	testRecordName := generateTestRecord(p, testDomain, true, testRecordTXT)
	registerRecordMockResponder(testDomain, testRecordName)

	rrSets := RRsets{}
	rrSetName := makeDomainCanonical(testRecordName)
	rrSets.Sets = []RRset{{Name: &rrSetName, Type: RRTypePtr(RRTypeTXT),
		ChangeType: ChangeTypePtr(ChangeTypeDelete)}}

	if err := p.Records.Patch(context.Background(), testDomain, &rrSets); err != nil {
		t.Errorf("%s", err)
	}
}
