package powerdns

import (
	"strings"
)

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

type RRset struct {
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	TTL        int       `json:"ttl"`
	ChangeType string    `json:"changetype"`
	Records    []Record  `json:"records"`
	Comments   []Comment `json:"comments"`
}

type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
	SetPTR   bool   `json:"set-ptr"`
}

type Comment struct {
	Content    string `json:"content"`
	Account    string `json:"account"`
	ModifiedAt int    `json:"modified_at"`
}

type RRsets struct {
	Sets []RRset `json:"rrsets"`
}

type NotifyResult struct {
	Result string `json:"result"`
}

func (p *PowerDNS) GetZones() ([]Zone, error) {
	zones := make([]Zone, 0)
	error := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost+"/zones").Receive(&zones, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	return zones, err
}

func (p *PowerDNS) GetZone(domain string) (*Zone, error) {
	zone := &Zone{}
	error := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost+"/zones/"+strings.TrimRight(domain, ".")).Receive(zone, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return &Zone{}, error
	}

	zone.PowerDNSHandle = p
	return zone, err
}

func (z *Zone) Notify() (*NotifyResult, error) {
	notifyResult := &NotifyResult{}
	error := new(Error)
	serversSling := z.PowerDNSHandle.makeSling()
	resp, err := serversSling.New().Put(strings.TrimRight(z.URL, ".")+"/notify").Receive(notifyResult, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return &NotifyResult{}, error
	}

	return notifyResult, err
}

func (z *Zone) AddRecord(name string, recordType string, ttl int, content []string) error {
	return z.ChangeRecord(name, recordType, ttl, content)
}

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

	error := new(Error)
	zone := new(Zone)

	zonesSling := z.PowerDNSHandle.makeSling()
	resp, err := zonesSling.New().Patch(strings.TrimRight(z.URL, ".")).BodyJSON(payload).Receive(zone, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return error
	}

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}

	return err
}
