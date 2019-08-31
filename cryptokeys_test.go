package powerdns

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func registerCryptokeyMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain+"/cryptokeys/0",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				cryptokeyMock := Cryptokey{
					Type:       "Cryptokey",
					ID:         0,
					KeyType:    "zsk",
					Active:     true,
					DNSkey:     "256 3 8 thisIsTheKey",
					Privatekey: "Private-key-format: v1.2\nAlgorithm: 8 (RSASHA256)\nModulus: foo\nPublicExponent: foo\nPrivateExponent: foo\nPrime1: foo\nPrime2: foo\nExponent1: foo\nExponent2: foo\nCoefficient: foo\n",
					Algorithm:  "RSASHA256",
					Bits:       1024,
				}
				return httpmock.NewJsonResponse(200, cryptokeyMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
}

func TestGetCryptokeys(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain+"/cryptokeys",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				cryptokeysMock := []Cryptokey{
					{
						Type:      "Cryptokey",
						ID:        11,
						KeyType:   "zsk",
						Active:    true,
						DNSkey:    "256 3 8 thisIsTheKey",
						Algorithm: "RSASHA256",
						Bits:      1024,
					},
					{
						Type:    "Cryptokey",
						ID:      10,
						KeyType: "lsk",
						Active:  true,
						DNSkey:  "257 3 8 thisIsTheKey",
						DS: []string{
							"997 8 1 foo",
							"997 8 2 foo",
							"997 8 4 foo",
						},
						Algorithm: "RSASHA256",
						Bits:      2048,
					},
				}
				return httpmock.NewJsonResponse(200, cryptokeysMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()

	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	cryptokeys, err := z.GetCryptokeys()
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(cryptokeys) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestGetCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)
	registerCryptokeyMockResponder(testDomain)

	p := initialisePowerDNSTestClient()

	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	cryptokey, err := z.GetCryptokey(0)
	if err != nil {
		t.Errorf("%s", err)
	}

	if cryptokey.Algorithm != "RSASHA256" {
		t.Error("Received cryptokey algorithm is wrong")
	}
}

func TestToggleCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)
	registerCryptokeyMockResponder(testDomain)
	httpmock.RegisterResponder("PUT", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain+"/cryptokeys/0",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(204, ""), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()

	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	c, err := z.GetCryptokey(0)
	if err != nil {
		t.Errorf("%s", err)
	}

	if c.ToggleCryptokey() != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)
	registerCryptokeyMockResponder(testDomain)
	httpmock.RegisterResponder("DELETE", "http://localhost:8080/api/v1/servers/localhost/zones/"+testDomain+"/cryptokeys/0",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(204, ""), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	p := initialisePowerDNSTestClient()

	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	c, err := z.GetCryptokey(0)
	if err != nil {
		t.Errorf("%s", err)
	}

	if c.DeleteCryptokey() != nil {
		t.Errorf("%s", err)
	}
}
