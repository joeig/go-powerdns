package powerdns

import (
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"math/rand"
	"net/http"
	"time"
)

const (
	testBaseURL string = "http://localhost:8080"
	testVhost   string = "localhost"
	testAPIKey  string = "apipw"
)

func generateTestAPIURL() string {
	return fmt.Sprintf("%s/api/v1", testBaseURL)
}

func generateTestAPIVhostURL() string {
	return fmt.Sprintf("%s/servers/%s", generateTestAPIURL(), testVhost)
}

func initialisePowerDNSTestClient() *Client {
	return NewClient(testBaseURL, testVhost, map[string]string{"X-API-Key": testAPIKey}, nil)
}

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

func registerZoneMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", generateTestAPIVhostURL()+"/zones/"+testDomain,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				zoneMock := Zone{
					ID:   String(fixDomainSuffix(testDomain)),
					Name: String(fixDomainSuffix(testDomain)),
					URL:  String("/api/v1/servers/" + testVhost + "/zones/" + fixDomainSuffix(testDomain)),
					Kind: ZoneKindPtr(NativeZoneKind),
					RRsets: []RRset{
						{
							Name: String(fixDomainSuffix(testDomain)),
							Type: String("SOA"),
							TTL:  Uint32(3600),
							Records: []Record{
								{
									Content: String("a.misconfigured.powerdns.server. hostmaster." + fixDomainSuffix(testDomain) + " 1337 10800 3600 604800 3600"),
								},
							},
						},
					},
					Serial:         Uint32(1337),
					NotifiedSerial: Uint32(1337),
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
}

func registerCryptokeysMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", generateTestAPIVhostURL()+"/zones/"+testDomain+"/cryptokeys",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				cryptokeysMock := []Cryptokey{
					{
						Type:      String("Cryptokey"),
						ID:        Uint64(11),
						KeyType:   String("zsk"),
						Active:    Bool(true),
						DNSkey:    String("256 3 8 thisIsTheKey"),
						Algorithm: String("ECDSAP256SHA256"),
						Bits:      Uint64(1024),
					},
					{
						Type:    String("Cryptokey"),
						ID:      Uint64(10),
						KeyType: String("lsk"),
						Active:  Bool(true),
						DNSkey:  String("257 3 8 thisIsTheKey"),
						DS: []string{
							"997 8 1 foo",
							"997 8 2 foo",
							"997 8 4 foo",
						},
						Algorithm: String("ECDSAP256SHA256"),
						Bits:      Uint64(2048),
					},
				}
				return httpmock.NewJsonResponse(200, cryptokeysMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
}

func registerCryptokeyMockResponder(testDomain string, id uint64) {
	httpmock.RegisterResponder("GET", generateTestAPIVhostURL()+"/zones/"+testDomain+"/cryptokeys/"+cryptokeyIDToString(id),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				cryptokeyMock := Cryptokey{
					Type:       String("Cryptokey"),
					ID:         Uint64(0),
					KeyType:    String("zsk"),
					Active:     Bool(true),
					DNSkey:     String("256 3 8 thisIsTheKey"),
					Privatekey: String("Private-key-format: v1.2\nAlgorithm: 8 (ECDSAP256SHA256)\nModulus: foo\nPublicExponent: foo\nPrivateExponent: foo\nPrime1: foo\nPrime2: foo\nExponent1: foo\nExponent2: foo\nCoefficient: foo\n"),
					Algorithm:  String("ECDSAP256SHA256"),
					Bits:       Uint64(1024),
				}
				return httpmock.NewJsonResponse(200, cryptokeyMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
}
