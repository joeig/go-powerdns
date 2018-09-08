package powerdns

import "strings"

// Cryptokey structure with JSON API metadata
type Cryptokey struct {
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	KeyType    string   `json:"keytype"`
	Active     bool     `json:"active"`
	DNSkey     string   `json:"dnskey"`
	DS         []string `json:"ds"`
	Privatekey string   `json:"privatekey"`
	Algorithm  string   `json:"algorithm"`
	Bits       int      `json:"bits"`
	ZoneHandle *Zone
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
func (z *Zone) GetCryptokey(id string) (*Cryptokey, error) {
	cryptokey := new(Cryptokey)
	myError := new(Error)
	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys/"+id).Receive(cryptokey, myError)

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
	resp, err := serversSling.New().Put(strings.TrimRight(c.ZoneHandle.URL, ".")+"/cryptokeys/"+c.ID).Receive(cryptokey, myError)

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
	resp, err := serversSling.New().Delete(strings.TrimRight(c.ZoneHandle.URL, ".")+"/cryptokeys/"+c.ID).Receive(cryptokey, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return myError
	}

	return nil
}
