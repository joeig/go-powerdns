package powerdns

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

func TestConvertCryptokeyIDToString(t *testing.T) {
	if cryptokeyIDToString(1337) != "1337" {
		t.Error("Cryptokey ID to string conversion failed")
	}
}

func TestGetCryptokeys(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerZoneMockResponder(testDomain)
	registerCryptokeysMockResponder(testDomain)

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

	p := initialisePowerDNSTestClient()

	registerZoneMockResponder(testDomain)
	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := z.GetCryptokeys()
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID

	registerCryptokeyMockResponder(testDomain, id)
	cryptokey, err := z.GetCryptokey(id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if cryptokey.Algorithm != "ECDSAP256SHA256" {
		t.Error("Received cryptokey algorithm is wrong")
	}
}

func TestToggleCryptokey(t *testing.T) {
	testDomain := generateTestZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	registerZoneMockResponder(testDomain)
	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := z.GetCryptokeys()
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID

	registerCryptokeyMockResponder(testDomain, id)
	httpmock.RegisterResponder("PUT", generateTestAPIVhostURL()+"/zones/"+testDomain+"/cryptokeys/"+cryptokeyIDToString(id),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewStringResponse(204, ""), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	c, err := z.GetCryptokey(id)
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

	p := initialisePowerDNSTestClient()

	registerZoneMockResponder(testDomain)

	z, err := p.GetZone(testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := z.GetCryptokeys()
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID

	registerCryptokeyMockResponder(testDomain, id)
	httpmock.RegisterResponder("DELETE", generateTestAPIVhostURL()+"/zones/"+testDomain+"/cryptokeys/"+cryptokeyIDToString(id),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewStringResponse(204, ""), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	c, err := z.GetCryptokey(id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if c.DeleteCryptokey() != nil {
		t.Errorf("%s", err)
	}
}
