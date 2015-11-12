package powerdns

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/url"
	"strings"
)

type Error struct {
	Message string `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

type CombinedRecord struct {
	Name    string
	Type    string
	TTL     int
	Records []string
}

// Zone
type Zone struct {
	Id             string `json:"id"`
	Url            string `json:"url"`
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

// RR
type Record struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
	Disabled bool   `json:"disabled"`
	Content  string `json:"content"`
}

// RRset
type RRset struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	ChangeType string   `json:"changetype"`
	Records    []Record `json:"records"`
}

// RRsets
type RRsets struct {
	Sets []RRset `json:"rrsets"`
}

type PowerDNS struct {
	hostname string
	port     string
	server   string
	apikey   string
}

// New returns a new PowerDNS
func New() *PowerDNS {
	return &PowerDNS{
		hostname: "localhost",
		port:     "8081",
		server:   "localhost",
		apikey:   "apikey",
	}
}

func (p *PowerDNS) Hostname(hostname string) *PowerDNS {
	p.hostname = hostname
	return p
}

func (p *PowerDNS) ApiKey(apikey string) *PowerDNS {
	p.apikey = apikey
	return p
}

func (p *PowerDNS) Port(port string) *PowerDNS {
	p.port = port
	return p
}

func (p *PowerDNS) Server(server string) *PowerDNS {
	p.server = server
	return p
}

func (p *PowerDNS) AddRecord(name string, record_type string, ttl int, records []string) (*Zone, error) {

	zone, err := p.ChangeRecord(name, record_type, ttl, records, "UPSERT")
	
	return zone, err
}

func (p *PowerDNS) DeleteRecord(name string, record_type string, ttl int, records []string) (*Zone, error) {

	zone, err := p.ChangeRecord(name, record_type, ttl, records, "DELETE")
	
	return zone, err
}

func (p *PowerDNS) ChangeRecord(name string, record_type string, ttl int, records []string, action string) (*Zone, error) {

	Record := new(CombinedRecord)
	Record.Name = name
	Record.Type = record_type
	Record.TTL = ttl
	Record.Records = records

	zone, err := p.patchRRset(*Record, action)
	
	return zone, err
}

func (p *PowerDNS) patchRRset(record CombinedRecord, action string) (*Zone, error) {

    if strings.HasSuffix(record.Name, ".") {
        record.Name = strings.TrimSuffix(record.Name, ".");
    }

    Set := RRset{ Name: record.Name, Type: record.Type, ChangeType: "REPLACE" }

    if action == "DELETE" {
        Set.ChangeType = "DELETE"
    }

    var R Record

    for _, rec := range record.Records {
        R = Record{ Name: record.Name, Type: record.Type, Content: rec }
        Set.Records = append(Set.Records, R)
    }

    json := RRsets{}
    json.Sets = append(json.Sets, Set)

    error := new(Error)
    zone := new(Zone)

    resp, err := p.getSling().Patch("rancher").BodyJSON(json).Receive(zone, error);

    if err == nil && resp.StatusCode >= 400 {
    	error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
    	return nil, error
    }

    return zone, err
}

func (p *PowerDNS) getSling() (*sling.Sling) {

	Url := new(url.URL)
	Url.Host = p.hostname + ":" + p.port
	Url.Scheme = "http"
	Url.Path = "/servers/" + p.server + "/zones/"

	Sling := sling.New().Base(Url.String())

	Sling.Set("X-API-Key", p.apikey)

	return Sling
}

func (p *PowerDNS) GetRecords(zone_name string) ([]Record, error) {

	var records []Record

	zone := new(Zone)
	error := new(Error)

	Url := new(url.URL)
	Url.Host = p.hostname + ":" + p.port
	Url.Scheme = "http"
	Url.Path = "/servers/" + p.server + "/zones/"

	resp, err := p.getSling().Path(zone_name).Set("X-API-Key", p.apikey).Receive(zone, error)

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

func (p *PowerDNS) GetCombinedRecords(ApiKey string) ([]CombinedRecord, error) {
	var records []CombinedRecord
	var uniqueRecords []CombinedRecord

	//- Plain records from the zone
	Records, err := p.GetRecords(ApiKey)

	if err != nil {
		return records, err
	}

	//- Iterate through records to combine them by name and type
	for _, rec := range Records {
		record := CombinedRecord{Name: rec.Name, Type: rec.Type, TTL: rec.TTL}
		found := false
		for _, u_rec := range uniqueRecords {
			if u_rec.Name == rec.Name && u_rec.Type == rec.Type {
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
	for _, u_rec := range uniqueRecords {
		for _, rec := range Records {
			if u_rec.Name == rec.Name && u_rec.Type == rec.Type {
				u_rec.Records = append(u_rec.Records, rec.Content)
			}
		}
		records = append(records, u_rec)
	}

	return records, nil
}

func init() {

}
