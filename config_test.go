package powerdns

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
)

func TestListConfig(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", generateTestAPIVHostURL()+"/config",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") != testAPIKey {
				return httpmock.NewStringResponse(http.StatusUnauthorized, "Unauthorized"), nil
			}

			configMock := []ConfigSetting{
				{
					Name:  String("signing-threads"),
					Type:  String("ConfigSetting"),
					Value: String("3"),
				},
			}
			return httpmock.NewJsonResponse(http.StatusOK, configMock)
		},
	)

	p := initialisePowerDNSTestClient()
	config, err := p.Config.List()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(config) == 0 {
		t.Error("Received amount of config settings is 0")
	}
}

func TestListConfigError(t *testing.T) {
	p := initialisePowerDNSTestClient()
	p.Port = "x"
	if _, err := p.Config.List(); err == nil {
		t.Error("error is nil")
	}
}
