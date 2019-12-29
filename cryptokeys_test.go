package powerdns

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func registerCryptokeysMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+testDomain+"/cryptokeys",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

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
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeysMock)
		},
	)
}

func registerCryptokeyMockResponder(testDomain string, id uint64) {
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+testDomain+"/cryptokeys/"+cryptokeyIDToString(id),
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

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
			return httpmock.NewJsonResponse(http.StatusOK, cryptokeyMock)
		},
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/zones/%s/cryptokeys/%s", generateTestAPIVHostURL(), testDomain, cryptokeyIDToString(id)),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)
}

func TestConvertCryptokeyIDToString(t *testing.T) {
	if cryptokeyIDToString(1337) != "1337" {
		t.Error("Cryptokey ID to string conversion failed")
	}
}

func TestListCryptokeys(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerCryptokeysMockResponder(testDomain)

	p := initialisePowerDNSTestClient()

	cryptokeys, err := p.Cryptokeys.List(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(cryptokeys) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestListCryptokeysError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Cryptokeys.List(testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestGetCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := p.Cryptokeys.List(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID

	registerCryptokeyMockResponder(testDomain, *id)
	cryptokey, err := p.Cryptokeys.Get(testDomain, *id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if *cryptokey.Algorithm != "ECDSAP256SHA256" {
		t.Error("Received cryptokey algorithm is wrong")
	}
}

func TestGetCryptokeyError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Cryptokeys.Get(testDomain, uint64(0)); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := p.Cryptokeys.List(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID
	registerCryptokeyMockResponder(testDomain, *id)
	if p.Cryptokeys.Delete(testDomain, *id) != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteCryptokeyError(t *testing.T) {
	testDomain := generateTestZone(false)
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if err := p.Cryptokeys.Delete(testDomain, uint64(0)); err == nil {
		t.Error("error is nil")
	}
}
