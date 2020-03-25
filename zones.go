package powerdns

import (
	"fmt"
	"io/ioutil"
)

// ZonesService handles communication with the zones related methods of the Client API
type ZonesService service

// Zone structure with JSON API metadata
type Zone struct {
	ID               *string   `json:"id,omitempty"`
	Name             *string   `json:"name,omitempty"`
	Type             *ZoneType `json:"type,omitempty"`
	URL              *string   `json:"url,omitempty"`
	Kind             *ZoneKind `json:"kind,omitempty"`
	RRsets           []RRset   `json:"rrsets,omitempty"`
	Serial           *uint32   `json:"serial,omitempty"`
	NotifiedSerial   *uint32   `json:"notified_serial,omitempty"`
	Masters          []string  `json:"masters,omitempty"`
	DNSsec           *bool     `json:"dnssec,omitempty"`
	Nsec3Param       *string   `json:"nsec3param,omitempty"`
	Nsec3Narrow      *bool     `json:"nsec3narrow,omitempty"`
	Presigned        *bool     `json:"presigned,omitempty"`
	SOAEdit          *string   `json:"soa_edit,omitempty"`
	SOAEditAPI       *string   `json:"soa_edit_api,omitempty"`
	APIRectify       *bool     `json:"api_rectify,omitempty"`
	Zone             *string   `json:"zone,omitempty"`
	Account          *string   `json:"account,omitempty"`
	Nameservers      []string  `json:"nameservers,omitempty"`
	MasterTSIGKeyIDs []string  `json:"master_tsig_key_ids,omitempty"`
	SlaveTSIGKeyIDs  []string  `json:"slave_tsig_key_ids,omitempty"`
}

// NotifyResult structure with JSON API metadata
type NotifyResult struct {
	Result *string `json:"result,omitempty"`
}

// Export string type
type Export string

// ZoneType string type
type ZoneType string

// ZoneZoneType sets the zone's type to zone
const ZoneZoneType ZoneType = "Zone"

// ZoneTypePtr is a helper function that allocates a new ZoneType value to store v and returns a pointer to it.
func ZoneTypePtr(v ZoneType) *ZoneType {
	return &v
}

// ZoneKind string type
type ZoneKind string

// ZoneKindPtr is a helper function that allocates a new ZoneKind value to store v and returns a pointer to it.
func ZoneKindPtr(v ZoneKind) *ZoneKind {
	return &v
}

const (
	// NativeZoneKind sets the zone's kind to native
	NativeZoneKind ZoneKind = "Native"
	// MasterZoneKind sets the zone's kind to master
	MasterZoneKind ZoneKind = "Master"
	// SlaveZoneKind sets the zone's kind to slave
	SlaveZoneKind ZoneKind = "Slave"
)

// List retrieves a list of Zones
func (z *ZonesService) List() ([]Zone, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones", z.client.VHost), nil, nil)
	if err != nil {
		return nil, err
	}

	zones := make([]Zone, 0)
	_, err = z.client.do(req, &zones)
	return zones, err
}

// Get returns a certain Zone for a given domain
func (z *ZonesService) Get(domain string) (*Zone, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, trimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	zone := &Zone{}
	_, err = z.client.do(req, &zone)
	return zone, err
}

// AddNative creates a new native zone
func (z *ZonesService) AddNative(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit, soaEditApi string, apiRectify bool, nameservers []string) (*Zone, error) {
	zone := Zone{
		Name:        String(domain),
		Kind:        ZoneKindPtr(NativeZoneKind),
		DNSsec:      Bool(dnssec),
		Nsec3Param:  String(nsec3Param),
		Nsec3Narrow: Bool(nsec3Narrow),
		SOAEdit:     String(soaEdit),
		SOAEditAPI:  String(soaEditApi),
		APIRectify:  Bool(apiRectify),
		Nameservers: nameservers,
	}
	return z.postZone(&zone)
}

// AddMaster creates a new master zone
func (z *ZonesService) AddMaster(domain string, dnssec bool, nsec3Param string, nsec3Narrow bool, soaEdit, soaEditApi string, apiRectify bool, nameservers []string) (*Zone, error) {
	zone := Zone{
		Name:        String(domain),
		Kind:        ZoneKindPtr(MasterZoneKind),
		DNSsec:      Bool(dnssec),
		Nsec3Param:  String(nsec3Param),
		Nsec3Narrow: Bool(nsec3Narrow),
		SOAEdit:     String(soaEdit),
		SOAEditAPI:  String(soaEditApi),
		APIRectify:  Bool(apiRectify),
		Nameservers: nameservers,
	}
	return z.postZone(&zone)
}

// AddSlave creates a new slave zone
func (z *ZonesService) AddSlave(domain string, masters []string) (*Zone, error) {
	zone := Zone{
		Name:    String(domain),
		Kind:    ZoneKindPtr(SlaveZoneKind),
		Masters: masters,
	}
	return z.postZone(&zone)
}

// Add pre-created zone
func (z *ZonesService) Add(zone *Zone) (*Zone, error) {
	return z.postZone(zone)
}

func (z *ZonesService) postZone(zone *Zone) (*Zone, error) {
	zone.Name = String(makeDomainCanonical(*zone.Name))
	zone.Type = ZoneTypePtr(ZoneZoneType)

	req, err := z.client.newRequest("POST", fmt.Sprintf("servers/%s/zones", z.client.VHost), nil, zone)
	if err != nil {
		return nil, err
	}

	createdZone := new(Zone)
	_, err = z.client.do(req, &createdZone)
	return createdZone, err
}

// Change modifies an existing zone
func (z *ZonesService) Change(domain string, zone *Zone) error {
	zone.ID = nil
	zone.Name = nil
	zone.Type = nil
	zone.URL = nil

	req, err := z.client.newRequest("PUT", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, trimDomain(domain)), nil, zone)
	if err != nil {
		return err
	}

	_, err = z.client.do(req, nil)
	return err
}

// Delete removes a certain Zone for a given domain
func (z *ZonesService) Delete(domain string) error {
	req, err := z.client.newRequest("DELETE", fmt.Sprintf("servers/%s/zones/%s", z.client.VHost, trimDomain(domain)), nil, nil)
	if err != nil {
		return err
	}

	_, err = z.client.do(req, nil)
	return err
}

// Notify sends a DNS notify packet to all slaves
func (z *ZonesService) Notify(domain string) (*NotifyResult, error) {
	req, err := z.client.newRequest("PUT", fmt.Sprintf("servers/%s/zones/%s/notify", z.client.VHost, trimDomain(domain)), nil, nil)
	if err != nil {
		return nil, err
	}

	notifyResult := &NotifyResult{}
	_, err = z.client.do(req, notifyResult)
	return notifyResult, err
}

// Export returns a BIND-like Zone file
func (z *ZonesService) Export(domain string) (Export, error) {
	req, err := z.client.newRequest("GET", fmt.Sprintf("servers/%s/zones/%s/export", z.client.VHost, trimDomain(domain)), nil, nil)
	if err != nil {
		return "", err
	}

	resp, err := z.client.do(req, nil)
	if err != nil {
		return "", err
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return Export(bodyBytes), nil
}
