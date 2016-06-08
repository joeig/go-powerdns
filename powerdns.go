package powerdns

//Based off of github.com/waynz0r/powerdns

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/dghubble/sling"
)

// Error strct
type Error struct {
	Message string `json:"error"`
}

// Error Returns
func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

// CombinedRecord strct
type CombinedRecord struct {
	Name    string
	Type    string
	TTL     int
	Records []string
}

// Zone struct
type Zone struct {
	ID             string `json:"id"`
	URL            string `json:"url"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	DNSsec         bool   `json:"dnssec"`
	Serial         int    `json:"serial"`
	NotifiedSerial int    `json:"notified_serial"`
	LastCheck      int    `json:"last_check"`
	Records        []struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		TTL      int    `json:"ttl"`
		Priority int    `json:"priority"`
		Disabled bool   `json:"disabled"`
		Content  string `json:"content"`
	} `json:"records"`
}

// Record struct
type Record struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
	Disabled bool   `json:"disabled"`
	Content  string `json:"content"`
}

// RRset struct
type RRset struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	ChangeType string   `json:"changetype"`
	Records    []Record `json:"records"`
}

// RRsets struct
type RRsets struct {
	Sets []RRset `json:"rrsets"`
}

// PowerDNS struct
type PowerDNS struct {
	scheme   string
	hostname string
	port     string
	vhost    string
	domain   string
	apikey   string
}

// New returns a new PowerDNS
func New(baseURL string, vhost string, domain string, apikey string) *PowerDNS {
	if vhost == "" {
		vhost = "localhost"
	}

	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("%s is not a valid url: %v", baseURL, err)
	}
	hp := strings.Split(u.Host, ":")
	hostname := hp[0]
	var port string
	if len(hp) > 1 {
		port = hp[1]
	} else {
		if u.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	return &PowerDNS{
		scheme:   u.Scheme,
		hostname: hostname,
		port:     port,
		vhost:    vhost,
		domain:   domain,
		apikey:   apikey,
	}
}

// AddRecord ...
func (p *PowerDNS) AddRecord(name string, recordType string, ttl int, content []string) (*Zone, error) {

	zone, err := p.ChangeRecord(name, recordType, ttl, content, "UPSERT")

	return zone, err
}

// DeleteRecord ...
func (p *PowerDNS) DeleteRecord(name string, recordType string, ttl int, content []string) (*Zone, error) {

	zone, err := p.ChangeRecord(name, recordType, ttl, content, "DELETE")

	return zone, err
}

// ChangeRecord ...
func (p *PowerDNS) ChangeRecord(name string, recordType string, ttl int, content []string, action string) (*Zone, error) {

	Record := new(CombinedRecord)
	Record.Name = name
	Record.Type = recordType
	Record.TTL = ttl
	Record.Records = content

	zone, err := p.patchRRset(*Record, action)

	return zone, err
}

func (p *PowerDNS) patchRRset(record CombinedRecord, action string) (*Zone, error) {

	if strings.HasSuffix(record.Name, ".") {
		record.Name = strings.TrimSuffix(record.Name, ".")
	}

	Set := RRset{Name: record.Name, Type: record.Type, ChangeType: "REPLACE"}

	if action == "DELETE" {
		Set.ChangeType = "DELETE"
	}

	var R Record

	for _, rec := range record.Records {
		R = Record{Name: record.Name, Type: record.Type, TTL: record.TTL, Content: rec}
		Set.Records = append(Set.Records, R)
	}

	json := RRsets{}
	json.Sets = append(json.Sets, Set)

	error := new(Error)
	zone := new(Zone)

	resp, err := p.getSling().Patch(p.domain).BodyJSON(json).Receive(zone, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	return zone, err
}

func (p *PowerDNS) getSling() *sling.Sling {

	u := new(url.URL)
	u.Host = p.hostname + ":" + p.port
	u.Scheme = p.scheme
	u.Path = "/servers/" + p.vhost + "/zones/"

	Sling := sling.New().Base(u.String())

	Sling.Set("X-API-Key", p.apikey)

	return Sling
}

// GetRecords ...
func (p *PowerDNS) GetRecords() ([]Record, error) {

	var records []Record

	zone := new(Zone)
	error := new(Error)

	u := new(url.URL)
	u.Host = p.hostname + ":" + p.port
	u.Scheme = p.scheme
	u.Path = "/servers/" + p.vhost + "/zones/"

	resp, err := p.getSling().Path(p.domain).Set("X-API-Key", p.apikey).Receive(zone, error)

	if err != nil {
		return records, fmt.Errorf("PowerDNS API call has failed: %v", err)
	}

	if resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return records, error
	}

	for _, rec := range zone.Records {
		record := Record{Name: rec.Name, Type: rec.Type, TTL: rec.TTL, Priority: rec.Priority, Disabled: rec.Disabled, Content: rec.Content}
		records = append(records, record)
	}

	return records, err
}

// GetCombinedRecords ...
func (p *PowerDNS) GetCombinedRecords() ([]CombinedRecord, error) {
	var records []CombinedRecord
	var uniqueRecords []CombinedRecord

	//- Plain records from the zone
	Records, err := p.GetRecords()

	if err != nil {
		return records, err
	}

	//- Iterate through records to combine them by name and type
	for _, rec := range Records {
		record := CombinedRecord{Name: rec.Name, Type: rec.Type, TTL: rec.TTL}
		found := false
		for _, uRec := range uniqueRecords {
			if uRec.Name == rec.Name && uRec.Type == rec.Type {
				found = true
				continue
			}
		}

		//- append them only if missing
		if found == false {
			uniqueRecords = append(uniqueRecords, record)
		}
	}

	//- Get all values from the unique records
	for _, uRec := range uniqueRecords {
		for _, rec := range Records {
			if uRec.Name == rec.Name && uRec.Type == rec.Type {
				uRec.Records = append(uRec.Records, rec.Content)
			}
		}
		records = append(records, uRec)
	}

	return records, nil
}

func init() {

}
