package powerdns

import (
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestGetCryptokeys(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   "example.com.",
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
					Kind: "Native",
					RRsets: []RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster.example.com. 1337 10800 3600 604800 3600",
								},
							},
						},
					},
					Serial:         1337,
					NotifiedSerial: 1337,
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/cryptokeys",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				cryptokeysMock := []Cryptokey{
					{
						Type:      "Cryptokey",
						ID:        "11",
						KeyType:   "zsk",
						Active:    true,
						DNSkey:    "256 3 8 thisIsTheKey",
						Algorithm: "RSASHA256",
						Bits:      1024,
					},
					{
						Type:    "Cryptokey",
						ID:      "10",
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

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	cryptokeys, err := z.GetCryptokeys()
	if err != nil {
		t.Errorf("%s", err)
	}
	if len(cryptokeys) == 0 {
		t.Error("Received amount of statistics is 0")
	}
}

func TestGetCryptokey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   "example.com.",
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
					Kind: "Native",
					RRsets: []RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster.example.com. 1337 10800 3600 604800 3600",
								},
							},
						},
					},
					Serial:         1337,
					NotifiedSerial: 1337,
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/cryptokeys/11",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				cryptokeyMock := Cryptokey{
					Type:       "Cryptokey",
					ID:         "11",
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

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	cryptokey, err := z.GetCryptokey("11")
	if err != nil {
		t.Errorf("%s", err)
	}
	if cryptokey.Algorithm != "RSASHA256" {
		t.Error("Received cryptokey algorithm is wrong")
	}
}

func TestToggleCryptokey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   "example.com.",
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
					Kind: "Native",
					RRsets: []RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster.example.com. 1337 10800 3600 604800 3600",
								},
							},
						},
					},
					Serial:         1337,
					NotifiedSerial: 1337,
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/cryptokeys/11",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				cryptokeyMock := Cryptokey{
					Type:       "Cryptokey",
					ID:         "11",
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
	httpmock.RegisterResponder("PUT", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/cryptokeys/11",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(204, ""), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	c, err := z.GetCryptokey("11")
	if c.ToggleCryptokey() != nil {
		t.Errorf("%s", err)
	}
}

func TestDeleteCryptokey(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				zoneMock := Zone{
					ID:   "example.com.",
					Name: "example.com.",
					URL:  "/api/v1/servers/localhost/zones/example.com.",
					Kind: "Native",
					RRsets: []RRset{
						{
							Name: "example.com.",
							Type: "SOA",
							TTL:  3600,
							Records: []Record{
								{
									Content: "a.misconfigured.powerdns.server. hostmaster.example.com. 1337 10800 3600 604800 3600",
								},
							},
						},
					},
					Serial:         1337,
					NotifiedSerial: 1337,
				}
				return httpmock.NewJsonResponse(200, zoneMock)
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)
	httpmock.RegisterResponder("GET", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/cryptokeys/11",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				cryptokeyMock := Cryptokey{
					Type:       "Cryptokey",
					ID:         "11",
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
	httpmock.RegisterResponder("DELETE", "http://localhost:8080/api/v1/servers/localhost/zones/example.com/cryptokeys/11",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("X-Api-Key") == "apipw" {
				return httpmock.NewStringResponse(204, ""), nil
			}
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	headers := make(map[string]string)
	headers["X-API-Key"] = "apipw"
	p := NewClient("http://localhost:8080/", "localhost", headers, nil)
	z, err := p.GetZone("example.com")
	c, err := z.GetCryptokey("11")
	if c.DeleteCryptokey() != nil {
		t.Errorf("%s", err)
	}
}
