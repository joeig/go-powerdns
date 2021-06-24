package powerdns

import (
	"fmt"
)

// RecordsService handles communication with the records related methods of the Client API
type RecordsService service

// RRset structure with JSON API metadata
type RRset struct {
	Name       *string     `json:"name,omitempty"`
	Type       *RRType     `json:"type,omitempty"`
	TTL        *uint32     `json:"ttl,omitempty"`
	ChangeType *ChangeType `json:"changetype,omitempty"`
	Records    []Record    `json:"records"`
	Comments   []Comment   `json:"comments,omitempty"`
}

// Record structure with JSON API metadata
type Record struct {
	Content  *string `json:"content,omitempty"`
	Disabled *bool   `json:"disabled,omitempty"`
	SetPTR   *bool   `json:"set-ptr,omitempty"`
}

// Comment structure with JSON API metadata
type Comment struct {
	Content    *string `json:"content,omitempty"`
	Account    *string `json:"account,omitempty"`
	ModifiedAt *uint64 `json:"modified_at,omitempty"`
}

// RRsets structure with JSON API metadata
type RRsets struct {
	Sets []RRset `json:"rrsets,omitempty"`
}

// ChangeType represents a string-valued change type
type ChangeType string

// ChangeTypePtr is a helper function that allocates a new ChangeType value to store v and returns a pointer to it.
func ChangeTypePtr(v ChangeType) *ChangeType {
	return &v
}

const (
	// ChangeTypeReplace represents the REPLACE change type
	ChangeTypeReplace ChangeType = "REPLACE"
	// ChangeTypeDelete represents the DELETE change type
	ChangeTypeDelete ChangeType = "DELETE"
)

// RRType represents a string-valued resource record type
type RRType string

// RRTypePtr is a helper function that allocates a new RRType value to store v and returns a pointer to it.
func RRTypePtr(v RRType) *RRType {
	return &v
}

const (
	// RRTypeA represents the A resource record type
	RRTypeA RRType = "A"
	// RRTypeAAAA represents the AAAA resource record type
	RRTypeAAAA RRType = "AAAA"
	// RRTypeA6 represents the A6 resource record type
	RRTypeA6 RRType = "A6"
	// RRTypeAFSDB represents the AFSDB resource record type
	RRTypeAFSDB RRType = "AFSDB"
	// RRTypeALIAS represents the ALIAS resource record type
	RRTypeALIAS RRType = "ALIAS"
	// RRTypeDHCID represents the DHCID resource record type
	RRTypeDHCID RRType = "DHCID"
	// RRTypeDLV represents the DLV resource record type
	RRTypeDLV RRType = "DLV"
	// RRTypeCAA represents the CAA resource record type
	RRTypeCAA RRType = "CAA"
	// RRTypeCERT represents the CERT resource record type
	RRTypeCERT RRType = "CERT"
	// RRTypeCDNSKEY represents the CDNSKEY resource record type
	RRTypeCDNSKEY RRType = "CDNSKEY"
	// RRTypeCDS represents the CDS resource record type
	RRTypeCDS RRType = "CDS"
	// RRTypeCNAME represents the CNAME resource record type
	RRTypeCNAME RRType = "CNAME"
	// RRTypeDNSKEY represents the DNSKEY resource record type
	RRTypeDNSKEY RRType = "DNSKEY"
	// RRTypeDNAME represents the DNAME resource record type
	RRTypeDNAME RRType = "DNAME"
	// RRTypeDS represents the DS resource record type
	RRTypeDS RRType = "DS"
	// RRTypeEUI48 represents the EUI48 resource record type
	RRTypeEUI48 RRType = "EUI48"
	// RRTypeEUI64 represents the EUI64 resource record type
	RRTypeEUI64 RRType = "EUI64"
	// RRTypeHINFO represents the HINFO resource record type
	RRTypeHINFO RRType = "HINFO"
	// RRTypeIPSECKEY represents the IPSECKEY resource record type
	RRTypeIPSECKEY RRType = "IPSECKEY"
	// RRTypeKEY represents the KEY resource record type
	RRTypeKEY RRType = "KEY"
	// RRTypeKX represents the KX resource record type
	RRTypeKX RRType = "KX"
	// RRTypeLOC represents the LOC resource record type
	RRTypeLOC RRType = "LOC"
	// RRTypeLUA represents the LUA resource record type
	RRTypeLUA RRType = "LUA"
	// RRTypeMAILA represents the MAILA resource record type
	RRTypeMAILA RRType = "MAILA"
	// RRTypeMAILB represents the MAILB resource record type
	RRTypeMAILB RRType = "MAILB"
	// RRTypeMINFO represents the MINFO resource record type
	RRTypeMINFO RRType = "MINFO"
	// RRTypeMR represents the MR resource record type
	RRTypeMR RRType = "MR"
	// RRTypeMX represents the MX resource record type
	RRTypeMX RRType = "MX"
	// RRTypeNAPTR represents the NAPTR resource record type
	RRTypeNAPTR RRType = "NAPTR"
	// RRTypeNS represents the NS resource record type
	RRTypeNS RRType = "NS"
	// RRTypeNSEC represents the NSEC resource record type
	RRTypeNSEC RRType = "NSEC"
	// RRTypeNSEC3 represents the NSEC3 resource record type
	RRTypeNSEC3 RRType = "NSEC3"
	// RRTypeNSEC3PARAM represents the NSEC3PARAM resource record type
	RRTypeNSEC3PARAM RRType = "NSEC3PARAM"
	// RRTypeOPENPGPKEY represents the OPENPGPKEY resource record type
	RRTypeOPENPGPKEY RRType = "OPENPGPKEY"
	// RRTypePTR represents the PTR resource record type
	RRTypePTR RRType = "PTR"
	// RRTypeRKEY represents the RKEY resource record type
	RRTypeRKEY RRType = "RKEY"
	// RRTypeRP represents the RP resource record type
	RRTypeRP RRType = "RP"
	// RRTypeRRSIG represents the RRSIG resource record type
	RRTypeRRSIG RRType = "RRSIG"
	// RRTypeSIG represents the SIG resource record type
	RRTypeSIG RRType = "SIG"
	// RRTypeSOA represents the SOA resource record type
	RRTypeSOA RRType = "SOA"
	// RRTypeSPF represents the SPF resource record type
	RRTypeSPF RRType = "SPF"
	// RRTypeSSHFP represents the SSHFP resource record type
	RRTypeSSHFP RRType = "SSHFP"
	// RRTypeSRV represents the SRV resource record type
	RRTypeSRV RRType = "SRV"
	// RRTypeTKEY represents the TKEY resource record type
	RRTypeTKEY RRType = "TKEY"
	// RRTypeTSIG represents the TSIG resource record type
	RRTypeTSIG RRType = "TSIG"
	// RRTypeTLSA represents the TLSA resource record type
	RRTypeTLSA RRType = "TLSA"
	// RRTypeSMIMEA represents the SMIMEA resource record type
	RRTypeSMIMEA RRType = "SMIMEA"
	// RRTypeTXT represents the TXT resource record type
	RRTypeTXT RRType = "TXT"
	// RRTypeURI represents the URI resource record type
	RRTypeURI RRType = "URI"
	// RRTypeWKS represents the WKS resource record type
	RRTypeWKS RRType = "WKS"
)

// Add creates a new resource record
func (r *RecordsService) Add(domain string, name string, recordType RRType, ttl uint32, content []string) error {
	return r.Change(domain, name, recordType, ttl, content)
}

// Change replaces an existing resource record
func (r *RecordsService) Change(domain string, name string, recordType RRType, ttl uint32, content []string) error {
	rrset := new(RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.TTL = &ttl
	rrset.ChangeType = ChangeTypePtr(ChangeTypeReplace)
	rrset.Records = make([]Record, 0)

	for _, c := range content {
		r := Record{Content: String(c), Disabled: Bool(false), SetPTR: Bool(false)}
		rrset.Records = append(rrset.Records, r)
	}

	payload := r.prepareRRSet(rrset)
	return r.patchRRSet(domain, payload)
}

// Delete removes an existing resource record
func (r *RecordsService) Delete(domain string, name string, recordType RRType) error {
	rrset := new(RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.ChangeType = ChangeTypePtr(ChangeTypeDelete)

	payload := r.prepareRRSet(rrset)
	return r.patchRRSet(domain, payload)
}

// Patch method makes patch of already prepared rrsets
func (r *RecordsService) Patch(domain string, rrSets *RRsets) error {
	for i := range rrSets.Sets {
		fixRRSet(&rrSets.Sets[i])
	}
	return r.patchRRSet(domain, rrSets)
}

func canonicalResourceRecordValues(records []Record) {
	for i := range records {
		records[i].Content = String(makeDomainCanonical(*records[i].Content))
	}
}

func fixRRSet(rrset *RRset) {
	if *rrset.Type != RRTypeCNAME && *rrset.Type != RRTypeMX {
		return
	}
	canonicalResourceRecordValues(rrset.Records)
}

func (r *RecordsService) prepareRRSet(rrSet *RRset) *RRsets {
	rrSet.Name = String(makeDomainCanonical(*rrSet.Name))

	fixRRSet(rrSet)

	payload := RRsets{}
	payload.Sets = append(payload.Sets, *rrSet)
	return &payload
}

func (r *RecordsService) patchRRSet(domain string, rrSets *RRsets) error {

	req, err := r.client.newRequest("PATCH", fmt.Sprintf("servers/%s/zones/%s", r.client.VHost, trimDomain(domain)), nil, &rrSets)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}
