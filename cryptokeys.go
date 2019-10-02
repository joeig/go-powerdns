package powerdns

import (
	"strconv"
	"strings"
)

// Cryptokey structure with JSON API metadata
type Cryptokey struct {
	Type       string   `json:"type,omitempty"`
	ID         uint64   `json:"id,omitempty"`
	KeyType    string   `json:"keytype,omitempty"`
	Active     bool     `json:"active,omitempty"`
	DNSkey     string   `json:"dnskey,omitempty"`
	DS         []string `json:"ds,omitempty"`
	Privatekey string   `json:"privatekey,omitempty"`
	Algorithm  string   `json:"algorithm,omitempty"`
	Bits       uint64   `json:"bits,omitempty"`
	ZoneHandle *Zone    `json:"-"`
}

func cryptokeyIDToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}

// GetCryptokeys retrieves a list of Cryptokeys that belong to a Zone
func (z *Zone) GetCryptokeys() ([]Cryptokey, error) {
	cryptokeys := make([]Cryptokey, 0)
	myError := new(Error)

	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys").Receive(&cryptokeys, myError)
	if err := handleAPIClientError(resp, &err, myError); err != nil {
		return nil, err
	}

	for i := range cryptokeys {
		cryptokeys[i].ZoneHandle = z
	}

	return cryptokeys, err
}

// GetCryptokey returns a certain Cryptokey instance of a given Zone
func (z *Zone) GetCryptokey(id uint64) (*Cryptokey, error) {
	cryptokey := new(Cryptokey)
	myError := new(Error)

	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys/"+cryptokeyIDToString(id)).Receive(cryptokey, myError)
	if err := handleAPIClientError(resp, &err, myError); err != nil {
		return nil, err
	}

	cryptokey.ZoneHandle = z
	return cryptokey, err
}

// ToggleCryptokey enables/disables a given Cryptokey
func (z *Zone) ToggleCryptokey(id uint64) error {
	cryptokey := new(Cryptokey)
	myError := new(Error)

	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys/"+cryptokeyIDToString(id)).Receive(cryptokey, myError)

	return handleAPIClientError(resp, &err, myError)
}

// ToggleCryptokey enables/disables a given Cryptokey
func (c *Cryptokey) ToggleCryptokey() error {
	return c.ZoneHandle.ToggleCryptokey(c.ID)
}

// DeleteCryptokey removes a given Cryptokey
func (z *Zone) DeleteCryptokey(id uint64) error {
	cryptokey := new(Cryptokey)
	myError := new(Error)

	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Delete(strings.TrimRight(z.URL, ".")+"/cryptokeys/"+cryptokeyIDToString(id)).Receive(cryptokey, myError)

	return handleAPIClientError(resp, &err, myError)
}

// DeleteCryptokey removes a given Cryptokey
func (c *Cryptokey) DeleteCryptokey() error {
	return c.ZoneHandle.DeleteCryptokey(c.ID)
}
