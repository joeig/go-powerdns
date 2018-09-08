package powerdns

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Zone structure with JSON API metadata
type Zone struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	URL            string   `json:"url"`
	Kind           string   `json:"kind"`
	RRsets         []RRset  `json:"rrsets"`
	Serial         int      `json:"serial"`
	NotifiedSerial int      `json:"notified_serial"`
	Masters        []string `json:"masters"`
	DNSsec         bool     `json:"dnssec"`
	Nsec3Param     string   `json:"nsec3param"`
	Nsec3Narrow    bool     `json:"nsec3narrow"`
	Presigned      bool     `json:"presigned"`
	SOAEdit        string   `json:"soa_edit"`
	SOAEditAPI     string   `json:"soa_edit_api"`
	APIRectify     bool     `json:"api_rectify"`
	Zone           string   `json:"zone"`
	Account        string   `json:"account"`
	Nameservers    []string `json:"nameservers"`
	PowerDNSHandle *PowerDNS
}

// RRset structure with JSON API metadata
type RRset struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	TTL        int       `json:"ttl"`
	ChangeType string    `json:"changetype"`
	Records    []Record  `json:"records"`
	Comments   []Comment `json:"comments"`
}

// Record structure with JSON API metadata
type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
	SetPTR   bool   `json:"set-ptr"`
}

// Comment structure with JSON API metadata
type Comment struct {
	Content    string `json:"content"`
	Account    string `json:"account"`
	ModifiedAt int    `json:"modified_at"`
}

// RRsets structure with JSON API metadata
type RRsets struct {
	Sets []RRset `json:"rrsets"`
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

// AddRecord creates a new resource record
func (z *Zone) AddRecord(name string, recordType string, ttl int, content []string) error {
	return z.ChangeRecord(name, recordType, ttl, content)
}

// ChangeRecord replaces an existing resource record
func (z *Zone) ChangeRecord(name string, recordType string, ttl int, content []string) error {
	rrset := new(RRset)
	rrset.Name = name
	rrset.Type = recordType
	rrset.TTL = ttl
	rrset.ChangeType = "REPLACE"

	for _, c := range content {
		r := Record{Content: c, Disabled: false, SetPTR: false}
		rrset.Records = append(rrset.Records, r)
	}

	return z.patchRRset(*rrset)
}

// DeleteRecord removes an existing resource record
func (z *Zone) DeleteRecord(name string, recordType string) error {
	rrset := new(RRset)
	rrset.Name = name
	rrset.Type = recordType
	rrset.ChangeType = "DELETE"

	return z.patchRRset(*rrset)
}

func (z *Zone) patchRRset(rrset RRset) error {
	if !strings.HasSuffix(rrset.Name, ".") {
		rrset.Name += "."
	}

	payload := RRsets{}
	payload.Sets = append(payload.Sets, rrset)

	myError := new(Error)
	zone := new(Zone)

	zonesSling := z.PowerDNSHandle.makeSling()
	resp, err := zonesSling.New().Patch(strings.TrimRight(z.URL, ".")).BodyJSON(payload).Receive(zone, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return myError
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	return err
}

// Export returns a BIND-like Zone file
func (z *Zone) Export() (Export, error) {
	myError := new(Error)
	exportSling := z.PowerDNSHandle.makeSling()
	req, err := exportSling.New().Get(strings.TrimRight(z.URL, ".") + "/export").Request()
	resp, err := http.DefaultClient.Do(req)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return "", myError
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	return Export(bodyBytes), nil
}
