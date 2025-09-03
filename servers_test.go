package powerdns

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

func registerServersMockResponder() {
	httpmock.RegisterResponder("GET", generateTestAPIURL()+"/servers",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			serversMock := []Server{
				{
					Type:       String("Server"),
					ID:         String(testVHost),
					DaemonType: String("authoritative"),
					Version:    String("4.1.2"),
					URL:        String("/api/v1/servers/" + testVHost),
					ConfigURL:  String("/api/v1/servers/" + testVHost + "/config{/config_setting}"),
					ZonesURL:   String("/api/v1/servers/" + testVHost + "/zones{/zone}"),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, serversMock)
		},
	)

	httpmock.RegisterResponder("GET", generateTestAPIVHostURL(),
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			serverMock := Server{
				Type:       String("Server"),
				ID:         String(testVHost),
				DaemonType: String("authoritative"),
				Version:    String("4.1.2"),
				URL:        String("/api/v1/servers/" + testVHost),
				ConfigURL:  String("/api/v1/servers/" + testVHost + "/config{/config_setting}"),
				ZonesURL:   String("/api/v1/servers/" + testVHost + "/zones{/zone}"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, serverMock)
		},
	)
}

func registerCacheFlushMockResponder(testDomain string) {
	httpmock.RegisterResponder("PUT", fmt.Sprintf("%s/cache/flush", generateTestAPIVHostURL()),
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.URL.Query().Get("domain") != makeDomainCanonical(testDomain) {
				return httpmock.NewStringResponse(http.StatusUnprocessableEntity, "Unprocessable Eneity"), nil
			}

			cacheFlushResultMock := CacheFlushResult{
				Count:  Uint32(1),
				Result: String("foo"),
			}
			return httpmock.NewJsonResponse(http.StatusOK, cacheFlushResultMock)
		},
	)
}

func TestListServers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerServersMockResponder()

	p := initialisePowerDNSTestClient()
	servers, err := p.Servers.List(context.Background())
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(servers) == 0 {
		t.Error("Received amount of servers is 0")
	}
}

func TestListServersError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Servers.List(context.Background()); err == nil {
		t.Error("error is nil")
	}
}

func TestGetServer(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerServersMockResponder()

	p := initialisePowerDNSTestClient()
	server, err := p.Servers.Get(context.Background(), testVHost)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *server.ID != testVHost {
		t.Error("Received no server")
	}
}

func TestGetServerError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Servers.Get(context.Background(), testVHost); err == nil {
		t.Error("error is nil")
	}
}

func TestCacheFlush(t *testing.T) {
	testDomain := generateNativeZone(true)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	registerCacheFlushMockResponder(testDomain)

	p := initialisePowerDNSTestClient()
	cacheFlushResult, err := p.Servers.CacheFlush(context.Background(), testVHost, testDomain)
	if err != nil {
		t.Errorf("%s", err)
	}
	if *cacheFlushResult.Count != 1 {
		t.Error("Received cache flush result is invalid")
	}
}

func TestCacheFlushResultError(t *testing.T) {
	testDomain := generateNativeZone(false)
	p := initialisePowerDNSTestClient()
	p.BaseURL = "://"
	if _, err := p.Servers.CacheFlush(context.Background(), testVHost, testDomain); err == nil {
		t.Error("error is nil")
	}
}
