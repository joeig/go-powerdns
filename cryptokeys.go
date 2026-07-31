package powerdns

import (
	"context"
	"net/http"
	"path"
	"strconv"
)

// CryptokeysService handles communication with the cryptokeys related methods of the Client API
type CryptokeysService service

// Cryptokey structure with JSON API metadata
type Cryptokey struct {
	Type       *string  `json:"type,omitempty"`
	ID         *uint64  `json:"id,omitempty"`
	KeyType    *string  `json:"keytype,omitempty"`
	Active     *bool    `json:"active,omitempty"`
	Published  *bool    `json:"published,omitempty"`
	DNSkey     *string  `json:"dnskey,omitempty"`
	DS         []string `json:"ds,omitempty"`
	CDS        []string `json:"cds,omitempty"`
	Privatekey *string  `json:"privatekey,omitempty"`
	Algorithm  *string  `json:"algorithm,omitempty"`
	Bits       *uint64  `json:"bits,omitempty"`
}

func cryptokeyIDToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// List retrieves a list of Cryptokeys that belong to a Zone
func (c *CryptokeysService) List(ctx context.Context, domain string) ([]Cryptokey, error) {
	req, err := c.client.newRequest(ctx, http.MethodGet, path.Join("servers", c.client.VHost, "zones", makeDomainCanonical(domain), "cryptokeys"), nil, nil)
	if err != nil {
		return nil, err
	}

	cryptokeys := make([]Cryptokey, 0)
	_, err = c.client.do(req, &cryptokeys)
	return cryptokeys, err
}

// Get returns a certain Cryptokey instance of a given Zone
func (c *CryptokeysService) Get(ctx context.Context, domain string, id uint64) (*Cryptokey, error) {
	req, err := c.client.newRequest(ctx, http.MethodGet, path.Join("servers", c.client.VHost, "zones", makeDomainCanonical(domain), "cryptokeys", cryptokeyIDToString(id)), nil, nil)
	if err != nil {
		return nil, err
	}

	cryptokey := new(Cryptokey)
	_, err = c.client.do(req, &cryptokey)
	return cryptokey, err
}

// Create adds a Cryptokey to a Zone. When cryptokey.Privatekey, cryptokey.Bits and
// cryptokey.Algorithm are unset, PowerDNS generates a new key based on the server
// defaults; otherwise the supplied key material is imported. cryptokey.KeyType is
// required and must be one of "ksk", "zsk" or "csk".
func (c *CryptokeysService) Create(ctx context.Context, domain string, cryptokey Cryptokey) (*Cryptokey, error) {
	req, err := c.client.newRequest(ctx, http.MethodPost, path.Join("servers", c.client.VHost, "zones", makeDomainCanonical(domain), "cryptokeys"), nil, cryptokey)
	if err != nil {
		return nil, err
	}

	responseCryptokey := new(Cryptokey)
	_, err = c.client.do(req, &responseCryptokey)
	return responseCryptokey, err
}

// Change (de)activates or (un)publishes an existing Cryptokey. Only cryptokey.Active
// and cryptokey.Published are honoured by the API. It returns no content on success.
func (c *CryptokeysService) Change(ctx context.Context, domain string, id uint64, cryptokey Cryptokey) error {
	req, err := c.client.newRequest(ctx, http.MethodPut, path.Join("servers", c.client.VHost, "zones", makeDomainCanonical(domain), "cryptokeys", cryptokeyIDToString(id)), nil, cryptokey)
	if err != nil {
		return err
	}

	_, err = c.client.do(req, nil)
	return err
}

// Delete removes a given Cryptokey
func (c *CryptokeysService) Delete(ctx context.Context, domain string, id uint64) error {
	req, err := c.client.newRequest(ctx, http.MethodDelete, path.Join("servers", c.client.VHost, "zones", makeDomainCanonical(domain), "cryptokeys", cryptokeyIDToString(id)), nil, nil)
	if err != nil {
		return err
	}

	_, err = c.client.do(req, nil)
	return err
}
