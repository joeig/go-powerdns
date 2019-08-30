package powerdns

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Zone structure with JSON API metadata
type Zone struct {
	ID               string    `json:"id,omitempty"`
	Name             string    `json:"name,omitempty"`
	Type             ZoneType  `json:"type,omitempty"`
	URL              string    `json:"url,omitempty"`
	Kind             ZoneKind  `json:"kind,omitempty"`
	RRsets           []RRset   `json:"rrsets,omitempty"`
	Serial           int       `json:"serial,omitempty"`
	NotifiedSerial   int       `json:"notified_serial,omitempty"`
	Masters          []string  `json:"masters,omitempty"`
	DNSsec           bool      `json:"dnssec,omitempty"`
	Nsec3Param       string    `json:"nsec3param,omitempty"`
	Nsec3Narrow      bool      `json:"nsec3narrow,omitempty"`
	Presigned        bool      `json:"presigned,omitempty"`
	SOAEdit          string    `json:"soa_edit,omitempty"`
	SOAEditAPI       string    `json:"soa_edit_api,omitempty"`
	APIRectify       bool      `json:"api_rectify,omitempty"`
	Zone             string    `json:"zone,omitempty"`
	Account          string    `json:"account,omitempty"`
	Nameservers      []string  `json:"nameservers,omitempty"`
	MasterTSIGKeyIDs []string  `json:"master_tsig_key_ids,omitempty"`
	SlaveTSIGKeyIDs  []string  `json:"slave_tsig_key_ids,omitempty"`
	PowerDNSHandle   *PowerDNS `json:"-"`
}

// NotifyResult structure with JSON API metadata
type NotifyResult struct {
	Result string `json:"result,omitempty"`
}

// Export string type
type Export string

// ZoneType string type
type ZoneType string

// ZoneZoneType sets the zone's type to zone
const ZoneZoneType ZoneType = "Zone"

// ZoneKind string type
type ZoneKind string

const (
	// NativeZoneKind sets the zone's kind to native
	NativeZoneKind ZoneKind = "Native"
	// MasterZoneKind sets the zone's kind to master
	MasterZoneKind ZoneKind = "Master"
	// SlaveZoneKind sets the zone's kind to slave
	SlaveZoneKind ZoneKind = "Slave"
)

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

// AddNativeZone creates a new native zone
func (p *PowerDNS) AddNativeZone(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit string, soaEditApi string, apiRectify bool, nameservers []string) (*Zone, error) {
	zone := Zone{
		Name:        domain,
		Kind:        NativeZoneKind,
		DNSsec:      dnssec,
		Nsec3Param:  nsec3Param,
		Nsec3Narrow: nsec3Narrow,
		SOAEdit:     soaEdit,
		SOAEditAPI:  soaEditApi,
		APIRectify:  apiRectify,
		Nameservers: nameservers,
	}
	return p.postZone(&zone)
}

// AddMasterZone creates a new master zone
func (p *PowerDNS) AddMasterZone(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit string, soaEditApi string, apiRectify bool, nameservers []string) (*Zone, error) {
	zone := Zone{
		Name:        domain,
		Kind:        MasterZoneKind,
		DNSsec:      dnssec,
		Nsec3Param:  nsec3Param,
		Nsec3Narrow: nsec3Narrow,
		SOAEdit:     soaEdit,
		SOAEditAPI:  soaEditApi,
		APIRectify:  apiRectify,
		Nameservers: nameservers,
	}
	return p.postZone(&zone)
}

// AddSlaveZone creates a new slave zone
func (p *PowerDNS) AddSlaveZone(domain string, masters []string) (*Zone, error) {
	zone := Zone{
		Name:    domain,
		Kind:    SlaveZoneKind,
		Masters: masters,
	}
	return p.postZone(&zone)
}

func (p *PowerDNS) postZone(zone *Zone) (*Zone, error) {
	zone.Name = fixDomainSuffix(zone.Name)
	zone.Type = fixZoneType(zone.Type)

	myError := new(Error)
	createdZone := new(Zone)

	zonesSling := p.makeSling()
	resp, err := zonesSling.New().Post("servers/"+p.VHost+"/zones").BodyJSON(zone).Receive(createdZone, myError)

	createdZone.PowerDNSHandle = p

	if err != nil {
		return createdZone, err
	}

	switch code := resp.StatusCode; {
	case code == 201:
		return createdZone, nil
	default:
		return createdZone, err
	}
}

// ChangeZone modifies an existing zone
func (p *PowerDNS) ChangeZone(zone *Zone) error {
	if zone.Name == "" {
		return &Error{"Name attribute missing"}
	}

	adjustedZone := *zone
	adjustedZone.ID = ""
	adjustedZone.Name = ""
	adjustedZone.Type = fixZoneType(zone.Type)
	adjustedZone.URL = ""

	myError := new(Error)

	zoneSling := p.makeSling()
	resp, err := zoneSling.New().Put("servers/"+p.VHost+"/zones/"+strings.TrimRight(zone.Name, ".")).BodyJSON(adjustedZone).Receive(nil, myError)

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return myError
	}

	return nil
}

// DeleteZone removes a certain Zone for a given domain
func (p *PowerDNS) DeleteZone(domain string) error {
	myError := new(Error)
	zoneSling := p.makeSling()
	resp, err := zoneSling.New().Delete("servers/"+p.VHost+"/zones/"+strings.TrimRight(domain, ".")).Receive(nil, myError)

	if err != nil {
		return err
	}

	if resp.StatusCode != 204 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return myError
	}

	return nil
}

// DeleteZone removes a certain Zone for a given domain
func (z *Zone) DeleteZone() error {
	return z.PowerDNSHandle.DeleteZone(z.Name)
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

func fixDomainSuffix(domain string) string {
	if !strings.HasSuffix(domain, ".") {
		domain += "."
	}
	return domain
}

func fixZoneType(zoneType ZoneType) ZoneType {
	if zoneType == "" {
		return ZoneZoneType
	}
	return zoneType
}
