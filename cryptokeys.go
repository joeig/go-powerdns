package powerdns

import (
	"fmt"
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
	DNSkey     *string  `json:"dnskey,omitempty"`
	DS         []string `json:"ds,omitempty"`
	Privatekey *string  `json:"privatekey,omitempty"`
	Algorithm  *string  `json:"algorithm,omitempty"`
	Bits       *uint64  `json:"bits,omitempty"`
}

func cryptokeyIDToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// List retrieves a list of Cryptokeys that belong to a Zone
func (c *CryptokeysService) List(domain string) ([]Cryptokey, error) {
	req, err := c.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s/cryptokeys", c.client.VHost, trimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	cryptokeys := make([]Cryptokey, 0)
	_, err = c.client.do(req, &cryptokeys)
	return cryptokeys, err
}

// Get returns a certain Cryptokey instance of a given Zone
func (c *CryptokeysService) Get(domain string, id uint64) (*Cryptokey, error) {
	req, err := c.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s/cryptokeys/%s", c.client.VHost, trimDomain(domain), cryptokeyIDToString(id)), nil, nil)
	if err != nil {
		return nil, err
	}

	cryptokey := new(Cryptokey)
	_, err = c.client.do(req, &cryptokey)
	return cryptokey, err
}

// Delete removes a given Cryptokey
func (c *CryptokeysService) Delete(domain string, id uint64) error {
	req, err := c.client.newRequest("DELETE", fmt.Sprintf("servers/%s/zones/%s/cryptokeys/%s", c.client.VHost, trimDomain(domain), cryptokeyIDToString(id)), nil, nil)
	if err != nil {
		return err
	}

	_, err = c.client.do(req, nil)
	return err
}
