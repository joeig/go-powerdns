package powerdns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
)

const (
	insecureKey = "ruTjBX2Jw/2BlE//5255fmKHaSRvLvp6p+YyDDAXThnBN/1Mz/VwMw+HQJVtkpDsAXvpPuNNZhucdKmhiOS4Tg=="
)

func generateTestTSIGKey(client *Client, name string, key string, autoAddTSIGKey bool) *TSIGKey {
	tsigKeyName := fmt.Sprintf("test-%d-%s", rand.Int(), name)
	newTSIGKey := TSIGKey{
		Name:      &tsigKeyName,
		ID:        String(tsigKeyName + "."),
		Algorithm: String("hmac-sha256"),
		Key:       String(key),
	}
	if autoAddTSIGKey && httpmock.Disabled() {
		tsigKey, err := client.TSIGKey.Create(context.Background(), tsigKeyName, "hmac-sha256", key)
		if err != nil {
			log.Printf("Error creating TSIG Key: %s: %v\n", name, err)
		} else {
			fmt.Printf("created TSIG Key: %s\n", *tsigKey.Name)
		}
		return tsigKey
	}
	if key == "" {
		newTSIGKey.Key = String(insecureKey)
	}
	return &newTSIGKey
}

func registerTSIGKeyMockResponder(tsigKeys *[]TSIGKey) {

	httpmock.RegisterResponder(http.MethodGet, generateTestAPIVHostURL()+"/tsigkeys",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}
			return httpmock.NewJsonResponse(http.StatusOK, tsigKeys)
		},
	)

	httpmock.RegisterResponder(http.MethodPost, generateTestAPIVHostURL()+"/tsigkeys",
		func(req *http.Request) (*http.Response, error) {
			if res := verifyAPIKey(req); res != nil {
				return res, nil
			}

			if req.Body == nil {
				log.Print("Request body is nil")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			var clientTsigkey TSIGKey
			if json.NewDecoder(req.Body).Decode(&clientTsigkey) != nil {
				log.Print("Cannot decode request body")
				return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
			}

			clientTsigkey.ID = String(*clientTsigkey.Name + ".")
			clientTsigkey.Type = String("TSIGKey")
			for _, serverTsigkey := range *tsigKeys {
				if *serverTsigkey.ID == *clientTsigkey.ID {
					return httpmock.NewBytesResponse(http.StatusConflict, []byte{}), nil
				}
			}

			return httpmock.NewJsonResponse(http.StatusCreated, clientTsigkey)
		},
	)

	for _, tsigkey := range *tsigKeys {

		httpmock.RegisterResponder(http.MethodGet, generateTestAPIVHostURL()+"/tsigkeys/"+*tsigkey.ID,
			func(req *http.Request) (*http.Response, error) {
				if res := verifyAPIKey(req); res != nil {
					return res, nil
				}

				return httpmock.NewJsonResponse(http.StatusOK, tsigkey)
			},
		)

		httpmock.RegisterResponder(http.MethodPut, generateTestAPIVHostURL()+"/tsigkeys/"+*tsigkey.ID,
			func(req *http.Request) (*http.Response, error) {
				if res := verifyAPIKey(req); res != nil {
					return res, nil
				}

				if req.Body == nil {
					log.Print("Request body is nil")
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}
				var clientTsigkey TSIGKey
				if json.NewDecoder(req.Body).Decode(&clientTsigkey) != nil {
					log.Print("Cannot decode request body")
					return httpmock.NewBytesResponse(http.StatusBadRequest, []byte{}), nil
				}

				return httpmock.NewJsonResponse(http.StatusOK, clientTsigkey)
			},
		)
		httpmock.RegisterResponder(http.MethodDelete, generateTestAPIVHostURL()+"/tsigkeys/"+*tsigkey.ID,
			func(req *http.Request) (*http.Response, error) {
				if res := verifyAPIKey(req); res != nil {
					return res, nil
				}

				return httpmock.NewBytesResponse(http.StatusNoContent, []byte{}), nil
			},
		)
	}
}

func TestCreateTSIGKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	testgen := generateTestTSIGKey(p, "testgen", "", false)
	testfull := generateTestTSIGKey(p, "testfull", "kjhOh/NLyHVwPxFwBrwYs043A99hFDycWwYi1Nj6R2bM3Rboh515yIbEIzzQx9Xod0W6nN8vnSAAvsysrgkOPw==", false)
	existingkey := generateTestTSIGKey(p, "existingkey", "", true)

	registerTSIGKeyMockResponder(&[]TSIGKey{
		*existingkey,
	})

	testCases := []struct {
		tsigkey     TSIGKey
		wantSuccess bool
	}{
		{
			tsigkey:     *testgen,
			wantSuccess: true,
		},
		{
			tsigkey:     *testfull,
			wantSuccess: true,
		},
		{
			tsigkey:     *existingkey,
			wantSuccess: false,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			respTSIGKey, err := p.TSIGKey.Create(context.Background(), *tc.tsigkey.Name, *tc.tsigkey.Algorithm, *tc.tsigkey.Key)
			if !tc.wantSuccess {
				if err == nil {
					t.Error("no error on duplicate key")
				}
				return
			}
			if err != nil {
				t.Error(err)
				return
			}

			if *respTSIGKey.Name != *tc.tsigkey.Name || *respTSIGKey.Algorithm != *tc.tsigkey.Algorithm {
				t.Errorf("input name or algorithm did not match with output: name: %s, %s; algorithm: %s, %s", *respTSIGKey.Name, *tc.tsigkey.Name, *respTSIGKey.Algorithm, *tc.tsigkey.Algorithm)
			}
		})
	}
}

func TestPatchTSIGKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	testPutTSIGKey := generateTestTSIGKey(p, "testput", "", true)

	registerTSIGKeyMockResponder(&[]TSIGKey{
		*testPutTSIGKey,
	})

	testPutTSIGKey.Key = String("1yWS55DxB2H40lded3/2IGnhbW6dCntvO+igEcP47n2ikD1EO03NDGKsKValitiqrtAmk41UbYVpREN23GYAdg==")

	_, err := p.TSIGKey.Change(context.Background(), *testPutTSIGKey.ID, *testPutTSIGKey)
	if err != nil {
		t.Error(err)
	}
}

func TestTSIGKeyErrorNewRequests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	p.Port = "x"
	registerTSIGKeyMockResponder(&[]TSIGKey{})

	t.Run("Test Get invalid request", func(t *testing.T) {
		_, err := p.TSIGKey.Get(context.Background(), "thiskeydoesnotexist.")
		if err == nil {
			t.Error("error is nil")
		}
	})
	t.Run("Test List invalid request", func(t *testing.T) {
		_, err := p.TSIGKey.List(context.Background())
		if err == nil {
			t.Error("error is nil")
		}
	})
	t.Run("Test Create invalid request", func(t *testing.T) {
		_, err := p.TSIGKey.Create(context.Background(), "test", "hmac-sha256", "")
		if err == nil {
			t.Error("error is nil")
		}
	})
	t.Run("Test Change invalid request", func(t *testing.T) {
		_, err := p.TSIGKey.Change(context.Background(), "test", TSIGKey{})
		if err == nil {
			t.Error("error is nil")
		}
	})
	t.Run("Test Delete invalid request", func(t *testing.T) {
		err := p.TSIGKey.Delete(context.Background(), "test")
		if err == nil {
			t.Error("error is nil")
		}
	})
}

func TestGetTSIGKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	getTSIGKey := generateTestTSIGKey(p, "getkey", "", true)
	registerTSIGKeyMockResponder(&[]TSIGKey{
		*getTSIGKey,
	})

	t.Run("Test Get", func(t *testing.T) {
		_, err := p.TSIGKey.Get(context.Background(), *getTSIGKey.ID)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test List", func(t *testing.T) {
		tsigKeyList, err := p.TSIGKey.List(context.Background())
		if err != nil {
			t.Error(err)
		}
		if len(tsigKeyList) == 0 {
			t.Error("expected at least one list item")
		}

	})
}

func TestDeleteTSIGKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()

	existingTSIGKey := generateTestTSIGKey(p, "deleteexisting", "", true)
	missingTSIGKey := generateTestTSIGKey(p, "deletemissing", "", false)

	registerTSIGKeyMockResponder(&[]TSIGKey{
		*existingTSIGKey,
	})

	t.Run("Remove existing TSIG Key", func(t *testing.T) {
		err := p.TSIGKey.Delete(context.Background(), *existingTSIGKey.ID)
		if err != nil {
			t.Errorf("expected successfull delete got error: %v", err)
			return
		}
	})

	t.Run("Remove non-existing TSIG Key", func(t *testing.T) {
		err := p.TSIGKey.Delete(context.Background(), *missingTSIGKey.ID)
		if err == nil {
			t.Errorf("expected err. but got nil")
		}
	})
}
