package powerdns

import (
	"fmt"
	"gopkg.in/jarcoal/httpmock.v1"
	"net/http"
	"testing"
)

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
	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/zones/%s/cryptokeys/%s", generateTestAPIVhostURL(), testDomain, cryptokeyIDToString(*id)),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)

	if p.Cryptokeys.Delete(testDomain, *id) != nil {
		t.Errorf("%s", err)
	}
}
