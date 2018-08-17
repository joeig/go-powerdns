package powerdns

import "strings"

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

func (z *Zone) GetCryptokeys() ([]Cryptokey, error) {
	cryptokeys := make([]Cryptokey, 0)
	error := new(Error)
	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys").Receive(&cryptokeys, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	for i := range cryptokeys {
		cryptokeys[i].ZoneHandle = z
	}

	return cryptokeys, err
}

func (z *Zone) GetCryptokey(id string) (*Cryptokey, error) {
	cryptokey := new(Cryptokey)
	error := new(Error)
	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Get(strings.TrimRight(z.URL, ".")+"/cryptokeys/"+id).Receive(cryptokey, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	cryptokey.ZoneHandle = z
	return cryptokey, err
}

func (c *Cryptokey) ToggleCryptokey() error {
	cryptokey := new(Cryptokey)
	error := new(Error)
	serversSling := c.ZoneHandle.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Put(strings.TrimRight(c.ZoneHandle.URL, ".")+"/cryptokeys/"+c.ID).Receive(cryptokey, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return error
	}

	return nil
}

func (c *Cryptokey) DeleteCryptokey() error {
	cryptokey := new(Cryptokey)
	error := new(Error)
	serversSling := c.ZoneHandle.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Delete(strings.TrimRight(c.ZoneHandle.URL, ".")+"/cryptokeys/"+c.ID).Receive(cryptokey, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return error
	}

	return nil
}
