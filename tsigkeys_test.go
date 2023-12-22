package powerdns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"testing"

	"github.com/jarcoal/httpmock"
)

var (
	serverKey1 = TSIGKey{
		ID:        String("testkeyonserver."),
		Name:      String("testkeyonserver"),
		Algorithm: String("hmac-sha256"),
		Key:       String("ruTjBX2Jw/2BlE//5255fmKHaSRvLvp6p+YyDDAXThnBN/1Mz/VwMw+HQJVtkpDsAXvpPuNNZhucdKmhiOS4Tg=="),
		Type:      String("TSIGKey"),
	}
	serverKey2 = TSIGKey{
		ID:        String("testkey2onserver."),
		Name:      String("testkey2onserver"),
		Algorithm: String("hmac-sha256"),
		Key:       String("ruTjBX2Jw/2BlE//5255fmKHaSRvLvp6p+YyDDAXThnBN/1Mz/VwMw+HQJVtkpDsAXvpPuNNZhucdKmhiOS4Tg=="),
		Type:      String("TSIGKey"),
	}
	serverKeyPatch = TSIGKey{
		ID:        String("testput."),
		Name:      String("testput"),
		Algorithm: String("hmac-sha256"),
		Key:       String("ruTjBX2Jw/2BlE//5255fmKHaSRvLvp6p+YyDDAXThnBN/1Mz/VwMw+HQJVtkpDsAXvpPuNNZhucdKmhiOS4Tg=="),
		Type:      String("TSIGKey"),
	}
)

func registerTSIGKeyMockReponder(tsigKeys *[]TSIGKey) {

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

				for _, serverTsigkey := range *tsigKeys {
					if *serverTsigkey.ID == path.Base(req.URL.Path) {
						continue
					}
					if *serverTsigkey.ID == *clientTsigkey.ID {
						return httpmock.NewBytesResponse(http.StatusConflict, []byte{}), nil
					}
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
	registerTSIGKeyMockReponder(&[]TSIGKey{
		serverKey1,
	})

	testCases := []struct {
		tsigkey     TSIGKey
		wantSuccess bool
	}{
		{
			TSIGKey{
				Name:      String("testgen"),
				Algorithm: String("hmac-sha256"),
				Key:       String(""),
			},
			true,
		},
		{
			TSIGKey{
				Name:      String("testfull"),
				Algorithm: String("hmac-sha256"),
				Key:       String("kjhOh/NLyHVwPxFwBrwYs043A99hFDycWwYi1Nj6R2bM3Rboh515yIbEIzzQx9Xod0W6nN8vnSAAvsysrgkOPw=="),
			},
			true,
		},
		{
			TSIGKey{
				Name:      String("testkeyonserver"),
				Algorithm: String("hmac-sha256"),
				Key:       String(""),
			},
			false,
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
	registerTSIGKeyMockReponder(&[]TSIGKey{
		serverKey1,
		serverKeyPatch,
	})

	testCases := []struct {
		tsigkey     TSIGKey
		wantSuccess bool
	}{
		{
			tsigkey: TSIGKey{
				ID:        String("testput."),
				Name:      String("testput"),
				Algorithm: String("hmac-sha256"),
				Key:       String("1yWS55DxB2H40lded3/2IGnhbW6dCntvO+igEcP47n2ikD1EO03NDGKsKValitiqrtAmk41UbYVpREN23GYAdg=="),
			},
			wantSuccess: false,
		},
		{
			tsigkey: TSIGKey{
				ID:        String("testkeyonserver."),
				Name:      String("testkeyonserver"),
				Algorithm: String("hmac-sha256"),
				Key:       String("1yWS55DxB2H40lded3/2IGnhbW6dCntvO+igEcP47n2ikD1EO03NDGKsKValitiqrtAmk41UbYVpREN23GYAdg=="),
				Type:      String("TSIGKey"),
			},
			wantSuccess: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase%d", i), func(t *testing.T) {
			_, err := p.TSIGKey.Change(context.Background(), "testkeyonserver.", tc.tsigkey)
			if !tc.wantSuccess {
				if err == nil {
					t.Error("expected error. but got nil")
				}
				return
			}
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestTSIGKeyErrorNewRequests(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	p.Port = "x"
	registerTSIGKeyMockReponder(&[]TSIGKey{})

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
	serverTSIGKeys := &[]TSIGKey{
		serverKey1,
		serverKey2,
	}

	registerTSIGKeyMockReponder(serverTSIGKeys)

	t.Run("Test Get", func(t *testing.T) {
		_, err := p.TSIGKey.Get(context.Background(), "testkeyonserver.")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Test List", func(t *testing.T) {
		tsigkeys, err := p.TSIGKey.List(context.Background())
		if err != nil {
			t.Error(err)
		}

		if len(*serverTSIGKeys) != len(tsigkeys) {
			t.Errorf("wrong amount of elements returned. expected %d, got %d", len(*serverTSIGKeys), len(tsigkeys))
			return
		}

	})
}

func TestDeleteTSIGKey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := initialisePowerDNSTestClient()
	registerTSIGKeyMockReponder(&[]TSIGKey{
		serverKey1,
		serverKey2,
	})

	t.Run("Remove existing TSIG Key", func(t *testing.T) {
		err := p.TSIGKey.Delete(context.Background(), *serverKey1.ID)
		if err != nil {
			t.Errorf("expected successfull delete got error: %v", err)
			return
		}
	})

	t.Run("Remove non-existing TSIG Key", func(t *testing.T) {
		err := p.TSIGKey.Delete(context.Background(), "doesnotexist.")
		if err == nil {
			t.Errorf("expected err. but got nil")
		}
	})
}
