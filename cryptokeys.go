package powerdns

import (
	"strconv"
	"strings"
)

// Cryptokey structure with JSON API metadata
type Cryptokey struct {
	Type       string   `json:"type,omitempty"`
	ID         int      `json:"id,omitempty"`
	KeyType    string   `json:"keytype,omitempty"`
	Active     bool     `json:"active,omitempty"`
	DNSkey     string   `json:"dnskey,omitempty"`
	DS         []string `json:"ds,omitempty"`
	Privatekey string   `json:"privatekey,omitempty"`
	Algorithm  string   `json:"algorithm,omitempty"`
	Bits       int      `json:"bits,omitempty"`
	ZoneHandle *Zone    `json:"-"`
}

// GetCryptokeys retrieves a list of Cryptokeys that belong to a Zone
func (z *Zone) GetCryptokeys() ([]Cryptokey, error) {
	cryptokeys := make([]Cryptokey, 0)
	myError := new(Error)
	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys").Receive(&cryptokeys, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return nil, myError
	}

	for i := range cryptokeys {
		cryptokeys[i].ZoneHandle = z
	}

	return cryptokeys, err
}

// GetCryptokey returns a certain Cryptokey instance of a given Zone
func (z *Zone) GetCryptokey(id int) (*Cryptokey, error) {
	cryptokey := new(Cryptokey)
	myError := new(Error)
	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys/"+strconv.Itoa(id)).Receive(cryptokey, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return nil, myError
	}

	cryptokey.ZoneHandle = z
	return cryptokey, err
}

// ToggleCryptokey enables/disables a given Cryptokey
func (c *Cryptokey) ToggleCryptokey() error {
	cryptokey := new(Cryptokey)
	myError := new(Error)
	serversSling := c.ZoneHandle.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Put(strings.TrimRight(c.ZoneHandle.URL, ".")+"/cryptokeys/"+strconv.Itoa(c.ID)).Receive(cryptokey, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return myError
	}

	return nil
}

// DeleteCryptokey removes a given Cryptokey
func (c *Cryptokey) DeleteCryptokey() error {
	cryptokey := new(Cryptokey)
	myError := new(Error)
	serversSling := c.ZoneHandle.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Delete(strings.TrimRight(c.ZoneHandle.URL, ".")+"/cryptokeys/"+strconv.Itoa(c.ID)).Receive(cryptokey, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return myError
	}

	return nil
}
