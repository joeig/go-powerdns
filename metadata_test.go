package powerdns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func generateTestMetadata(client *Client, domain string) string {
	if httpmock.Disabled() {
		metadata, err := client.Metadata.Create(context.Background(), domain, MetadataAllowAXFRFrom, []string{"192.168.0.1", "::1"})
		if err != nil {
			fmt.Println("Error creating test metadata")
			fmt.Printf("%v\n", err)
			fmt.Printf("%v\n", metadata)
		} else {
			fmt.Println("Created test metadata")
		}
	}

	return domain
}

func registerMetadataMockResponder(testDomain string) {
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/metadata",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			metadataKindsMock := []Metadata{
				{
					Kind:     MetadataKindPtr(MetadataAllowAXFRFrom),
					Metadata: []string{"127.0.0.1", "::1"},
				},
				{
					Kind:     MetadataKindPtr(MetadataTSIGAllowAXFR),
					Metadata: []string{"127.0.0.1", "::1"},
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, metadataKindsMock)
		},
	)

	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/metadata/ALLOW-AXFR-FROM",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			metadataMock := Metadata{
				Kind:     MetadataKindPtr(MetadataAllowAXFRFrom),
				Metadata: []string{"127.0.0.1", "::1"},
			}
			return httpmock.NewJsonResponse(http.StatusOK, metadataMock)
		},
	)

	httpmock.RegisterResponder("PUT", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/metadata/ALLOW-AXFR-FROM",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var metadata Metadata
			if json.NewDecoder(req.Body).Decode(&metadata) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewJsonResponse(http.StatusOK, metadata)
		},
	)

	httpmock.RegisterResponder("POST", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/metadata",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var metadata Metadata
			if json.NewDecoder(req.Body).Decode(&metadata) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewJsonResponse(http.StatusCreated, metadata)
		},
	)

	httpmock.RegisterResponder("DELETE", generateTestAPIVHostURL()+"/zones/"+makeDomainCanonical(testDomain)+"/metadata/ALLOW-AXFR-FROM",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body != nil {
				log.Print("Request body is not nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			return httpmock.NewBytesResponse(http.StatusNoContent, []byte{}), nil
		},
	)
}

func TestListMetadata(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerMetadataMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	metadataKinds, err := p.Metadata.List(context.Background(), testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(metadataKinds) == 0 {
		t.Error("Received amount of metadata kinds is 0")
	}
}

func TestListMetadataError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Metadata.List(context.Background(), testDomain); err == nil {
		t.Error("error is nil")
	}
}

func TestCreateMetadata(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerMetadataMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	createdMetadata, err := p.Metadata.Create(context.Background(), testDomain, MetadataAllowAXFRFrom, []string{"192.168.0.1", "::1"})
	if err != nil {
		t.Errorf("%s", err)
	}
	if *createdMetadata.Kind != "ALLOW-AXFR-FROM" {
		t.Error("Received wrong metadata kind")
	}
	if len(createdMetadata.Metadata) != 2 {
		t.Error("Received wrong number of metadata values")
	}
}

func TestCreateMetadataError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Metadata.Create(context.Background(), testDomain, MetadataAllowAXFRFrom, []string{"192.168.0.1"}); err == nil {
		t.Error("error is nil")
	}
}

func TestGetMetadata(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerMetadataMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	generateTestMetadata(p, testDomain)
	metadata, err := p.Metadata.Get(context.Background(), testDomain, "ALLOW-AXFR-FROM")
	if err != nil {
		t.Errorf("%s", err)
	}
	if *metadata.Kind != "ALLOW-AXFR-FROM" {
		t.Error("Received wrong metadata kind")
	}
	if len(metadata.Metadata) != 2 {
		t.Error("Received wrong number of metadata values")
	}
}

func TestGetMetadataError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Metadata.Get(context.Background(), testDomain, "ALLOW-AXFR-FROM"); err == nil {
		t.Error("error is nil")
	}
}

func TestSetMetadata(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerMetadataMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	metadata, err := p.Metadata.Set(context.Background(), testDomain, "ALLOW-AXFR-FROM", []string{"192.168.0.1", "::1"})
	if err != nil {
		t.Errorf("%s", err)
	}
	if *metadata.Kind != "ALLOW-AXFR-FROM" {
		t.Error("Received wrong metadata kind")
	}
}

func TestSetMetadataError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Metadata.Set(context.Background(), testDomain, "ALLOW-AXFR-FROM", []string{"192.168.0.1"}); err == nil {
		t.Error("error is nil")
	}
}

func TestDeleteMetadata(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerMetadataMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	err := p.Metadata.Delete(context.Background(), testDomain, "ALLOW-AXFR-FROM")
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteMetadataError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if err := p.Metadata.Delete(context.Background(), testDomain, "ALLOW-AXFR-FROM"); err == nil {
		t.Error("error is nil")
	}
}
