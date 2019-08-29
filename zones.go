package powerdns

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Zone structure with JSON API metadata
type Zone struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	URL            string    `json:"url"`
	Kind           string    `json:"kind"`
	RRsets         []RRset   `json:"rrsets"`
	Serial         int       `json:"serial"`
	NotifiedSerial int       `json:"notified_serial"`
	Masters        []string  `json:"masters"`
	DNSsec         bool      `json:"dnssec"`
	Nsec3Param     string    `json:"nsec3param"`
	Nsec3Narrow    bool      `json:"nsec3narrow"`
	Presigned      bool      `json:"presigned"`
	SOAEdit        string    `json:"soa_edit"`
	SOAEditAPI     string    `json:"soa_edit_api"`
	APIRectify     bool      `json:"api_rectify"`
	Zone           string    `json:"zone"`
	Account        string    `json:"account"`
	Nameservers    []string  `json:"nameservers"`
	PowerDNSHandle *PowerDNS `json:"-"`
}

// NotifyResult structure with JSON API metadata
type NotifyResult struct {
	Result string `json:"result"`
}

// Export string type
type Export string

// GetZones retrieves a list of Zones
func (p *PowerDNS) GetZones() ([]Zone, error) {
	zones := make([]Zone, 0)
	myError := new(Error)
	zonesSling := p.makeSling()
	resp, err := zonesSling.New().Get("servers/"+p.VHost+"/zones").Receive(&zones, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return nil, myError
	}

	for i := range zones {
		zones[i].PowerDNSHandle = p
	}

	return zones, err
}

// GetZone returns a certain Zone for a given domain
func (p *PowerDNS) GetZone(domain string) (*Zone, error) {
	zone := &Zone{}
	myError := new(Error)
	zoneSling := p.makeSling()
	resp, err := zoneSling.New().Get("servers/"+p.VHost+"/zones/"+strings.TrimRight(domain, ".")).Receive(zone, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return &Zone{}, myError
	}

	zone.PowerDNSHandle = p
	return zone, err
}

// Notify sends a DNS notify packet to all slaves
func (z *Zone) Notify() (*NotifyResult, error) {
	notifyResult := &NotifyResult{}
	myError := new(Error)
	notifySling := z.PowerDNSHandle.makeSling()
	resp, err := notifySling.New().Put(strings.TrimRight(z.URL, ".")+"/notify").Receive(notifyResult, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return &NotifyResult{}, myError
	}

	return notifyResult, err
}

// Export returns a BIND-like Zone file
func (z *Zone) Export() (Export, error) {
	myError := new(Error)
	exportSling := z.PowerDNSHandle.makeSling()
	req, _ := exportSling.New().Get(strings.TrimRight(z.URL, ".") + "/export").Request()
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return "", myError
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return Export(bodyBytes), nil
}
