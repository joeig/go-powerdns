package powerdns

import "strings"

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

func (p *PowerDNS) AddRecord(name string, recordType string, ttl int, content []string) (*Zone, error) {
	return p.ChangeRecord(name, recordType, ttl, content)
}

func (p *PowerDNS) ChangeRecord(name string, recordType string, ttl int, content []string) (*Zone, error) {
	rrset := new(RRset)
	rrset.Name = name
	rrset.Type = recordType
	rrset.TTL = ttl
	rrset.ChangeType = "REPLACE"

	for _, c := range content {
		r := Record{Content: c, Disabled: false, SetPTR: false}
		rrset.Records = append(rrset.Records, r)
	}

	zone, err := p.patchRRset(*rrset)

	return zone, err
}

func (p *PowerDNS) DeleteRecord(name string, recordType string) (*Zone, error) {
	rrset := new(RRset)
	rrset.Name = name
	rrset.Type = recordType
	rrset.ChangeType = "DELETE"

	zone, err := p.patchRRset(*rrset)

	return zone, err
}

func (p *PowerDNS) patchRRset(rrset RRset) (*Zone, error) {
	if !strings.HasSuffix(rrset.Name, ".") {
		rrset.Name += "."
	}

	payload := RRsets{}
	payload.Sets = append(payload.Sets, rrset)

	error := new(Error)
	zone := new(Zone)

	zonesSling := p.makeSling(p.vhost + "/zones/")
	resp, err := zonesSling.New().Patch(p.domain).BodyJSON(payload).Receive(zone, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	return zone, err
}
