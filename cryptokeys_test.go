package powerdns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerCryptokeysMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/cryptokeys",
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
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/cryptokeys/"+cryptokeyIDToString(id),
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

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("%s/zones/%s/cryptokeys/%s", generateTestAPIVHostURL(), makeDomainCanonical(testDomain), cryptokeyIDToString(id)),
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == testAPIKey {
				return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
		},
	)
}

func registerCreateCryptokeyMockResponder(testDomain string) {
	httpmock.RegisterResponder("POST", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/cryptokeys",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			var reqCryptokey Cryptokey
			if err := json.NewDecoder(req.Body).Decode(&reqCryptokey); err != nil {
				return httpmock.NewStringResponse(http.StatusBadRequest, "Bad Request"), nil
			}

			if reqCryptokey.KeyType == nil {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "keytype is required"), nil
			}

			responseCryptokey := Cryptokey{
				Type:      String("Cryptokey"),
				ID:        Uint64(12),
				KeyType:   reqCryptokey.KeyType,
				Active:    Bool(true),
				Published: Bool(true),
				DNSkey:    String("257 3 13 thisIsTheNewKey"),
				DS: []string{
					"997 13 2 foo",
				},
				CDS: []string{
					"997 13 2 foo",
				},
				Algorithm: String("ECDSAP256SHA256"),
				Bits:      Uint64(256),
			}
			return httpmock.NewJsonResponse(http.StatusCreated, responseCryptokey)
		},
	)
}

func registerChangeCryptokeyMockResponder(testDomain string, id uint64) {
	httpmock.RegisterResponder("PUT", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/cryptokeys/"+cryptokeyIDToString(id),
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			var reqCryptokey Cryptokey
			if err := json.NewDecoder(req.Body).Decode(&reqCryptokey); err != nil {
				return httpmock.NewStringResponse(http.StatusBadRequest, "Bad Request"), nil
			}

			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		},
	)
}

func TestConvertCryptokeyIDToString(t *testing.T) {
	if cryptokeyIDToString(1337) != "1337" {
		t.Error("Cryptokey ID to string conversion failed")
	}
}

func TestListCryptokeys(t *testing.T) {
	testDomain := generateNativeZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerCryptokeysMockResponder(testDomain)

	p := initialisePowerDNSTestClient()

	cryptokeys, err := p.Cryptokeys.List(context.Background(), testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	if len(cryptokeys) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestListCryptokeysError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Cryptokeys.List(context.Background(), testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestGetCryptokey(t *testing.T) {
	testDomain := generateNativeZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := p.Cryptokeys.List(context.Background(), testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID

	registerCryptokeyMockResponder(testDomain, *id)
	cryptokey, err := p.Cryptokeys.Get(context.Background(), testDomain, *id)
	if err != nil {
		t.Errorf("%s", err)
	}

	if *cryptokey.Algorithm != "ECDSAP256SHA256" {
		t.Error("Received cryptokey algorithm is wrong")
	}
}

func TestGetCryptokeyError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Cryptokeys.Get(context.Background(), testDomain, uint64(0)); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteCryptokey(t *testing.T) {
	testDomain := generateNativeZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := p.Cryptokeys.List(context.Background(), testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID
	registerCryptokeyMockResponder(testDomain, *id)
	if err = p.Cryptokeys.Delete(context.Background(), testDomain, *id); err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteCryptokeyError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if err := p.Cryptokeys.Delete(context.Background(), testDomain, uint64(0)); err == nil {
		t.Error("error is nil")
	}
}

func TestCreateCryptokey(t *testing.T) {
	testDomain := generateNativeZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	registerCreateCryptokeyMockResponder(testDomain)

	cryptokey, err := p.Cryptokeys.Create(context.Background(), testDomain, Cryptokey{
		KeyType:   String("ksk"),
		Active:    Bool(true),
		Published: Bool(true),
	})
	if err != nil {
		t.Fatalf("%s", err)
	}

	if cryptokey.ID == nil {
		t.Error("Received cryptokey ID is nil")
	}

	if cryptokey.KeyType == nil {
		t.Error("Received cryptokey keytype is nil")
	}

	if cryptokey.Active == nil || !*cryptokey.Active {
		t.Error("Received cryptokey is not active")
	}

	if cryptokey.Published == nil || !*cryptokey.Published {
		t.Error("Received cryptokey is not published")
	}

	// The mock returns fixed values that a live server does not reproduce
	// (the ID is server-assigned, keytype is derived from the key's role in
	// the zone, and CDS is only populated when CDS publication is enabled).
	if !httpmock.Disabled() {
		if *cryptokey.ID != 12 {
			t.Error("Received cryptokey ID is wrong")
		}

		if *cryptokey.KeyType != "ksk" {
			t.Error("Received cryptokey keytype is wrong")
		}

		if len(cryptokey.CDS) == 0 {
			t.Error("Received cryptokey CDS is empty")
		}
	}
}

func TestCreateCryptokeyError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Cryptokeys.Create(context.Background(), testDomain, Cryptokey{KeyType: String("ksk")}); err == nil {
		t.Error("error is nil")
	}
}

func TestChangeCryptokey(t *testing.T) {
	testDomain := generateNativeZone(true)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	registerCryptokeysMockResponder(testDomain)
	cryptokeys, err := p.Cryptokeys.List(context.Background(), testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}

	id := cryptokeys[0].ID
	registerChangeCryptokeyMockResponder(testDomain, *id)

	if err = p.Cryptokeys.Change(context.Background(), testDomain, *id, Cryptokey{Active: Bool(false)}); err != nil {
		t.Errorf("%s", err)
	}
}

func TestChangeCryptokeyError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if err := p.Cryptokeys.Change(context.Background(), testDomain, uint64(0), Cryptokey{Active: Bool(false)}); err == nil {
		t.Error("error is nil")
	}
}
