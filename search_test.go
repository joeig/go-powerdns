package powerdns

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func generateSearchTestZone() string {
	domain := fmt.Sprintf("test-%d.com", rand.Int())

	if httpmock.Disabled() {
		pdns := initialisePowerDNSTestClient()
		_, err := pdns.Zones.AddNative(context.Background(), domain, false, "", false, "", "", true, []string{"ns.foo.tld."})
		if err != nil {
			fmt.Printf("Error creating search test zone %s: %v\n", domain, err)
		} else {
			fmt.Printf("Created search test zone %s\n", domain)
		}
	}

	return domain
}

func registerSearchMockResponder() {
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/search-data",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			query := req.URL.Query().Get("q")
			if query == "" {
				return httpmock.NewStringResponse(http.StatusBadRequest, "Bad Request"), nil
			}

			searchResultsMock := []SearchResult{
				{
					Content:    String("192.0.2.1"),
					Disabled:   Bool(false),
					Name:       String("www.example.com."),
					ObjectType: String("record"),
					ZoneID:     String("example.com."),
					Zone:       String("example.com."),
					Type:       String("A"),
					TTL:        Uint32(3600),
				},
				{
					Content:    String(""),
					Disabled:   Bool(false),
					Name:       String("example.com."),
					ObjectType: String("zone"),
					ZoneID:     String("example.com."),
					Zone:       String("example.com."),
					Type:       nil,
					TTL:        nil,
				},
			}

			objectType := req.URL.Query().Get("object_type")
			if objectType == "zone" {
				searchResultsMock = []SearchResult{
					{
						Content:    String(""),
						Disabled:   Bool(false),
						Name:       String("example.com."),
						ObjectType: String("zone"),
						ZoneID:     String("example.com."),
						Zone:       String("example.com."),
						Type:       nil,
						TTL:        nil,
					},
				}
			}

			return httpmock.NewJsonResponse(http.StatusOK, searchResultsMock)
		},
	)
}

func TestSearch(t *testing.T) {
	testDomain := generateSearchTestZone()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerSearchMockResponder()

	p := initialisePowerDNSTestClient()
	searchQuery := "example*"
	if httpmock.Disabled() {
		searchQuery = testDomain[:len(testDomain)-4] + "*" // Use the generated domain pattern
	}
	results, err := p.Search.Data(context.Background(), searchQuery, 100, SearchObjectTypeAll)
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(results) == 0 {
		t.Error("Received amount of search results is 0")
	}
}

func TestSearchWithObjectType(t *testing.T) {
	testDomain := generateSearchTestZone()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerSearchMockResponder()

	p := initialisePowerDNSTestClient()
	searchQuery := "example*"
	if httpmock.Disabled() {
		searchQuery = testDomain[:len(testDomain)-4] + "*" // Use the generated domain pattern
	}
	results, err := p.Search.Data(context.Background(), searchQuery, 100, SearchObjectTypeZone)
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(results) != 1 {
		t.Errorf("Received amount of search results is not 1, got %d", len(results))
	}
	if *results[0].ObjectType != "zone" {
		t.Error("Received search result is not a zone")
	}
}

func TestSearchError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Search.Data(context.Background(), "example*", 100, SearchObjectTypeAll); err == nil {
		t.Error("error is nil")
	}
}
