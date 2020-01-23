package powerdns

import (
	"fmt"
)

// RecordsService handles communication with the records related methods of the Client API
type RecordsService service

// RRset structure with JSON API metadata
type RRset struct {
	Name       *string   `json:"name,omitempty"`
	Type       *string   `json:"type,omitempty"`
	TTL        *uint32   `json:"ttl,omitempty"`
	ChangeType *string   `json:"changetype,omitempty"`
	Records    []Record  `json:"records"`
	Comments   []Comment `json:"comments,omitempty"`
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

// Add creates a new resource record
func (r *RecordsService) Add(domain string, name string, recordType string, ttl uint32, content []string) error {
	return r.Change(domain, name, recordType, ttl, content)
}

// Change replaces an existing resource record
func (r *RecordsService) Change(domain string, name string, recordType string, ttl uint32, content []string) error {
	rrset := new(RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.TTL = &ttl
	rrset.ChangeType = String("REPLACE")
	rrset.Records = make([]Record, 0)

	for _, c := range content {
		r := Record{Content: String(c), Disabled: Bool(false), SetPTR: Bool(false)}
		rrset.Records = append(rrset.Records, r)
	}

	return r.patchRRset(domain, *rrset)
}

// Delete removes an existing resource record
func (r *RecordsService) Delete(domain string, name string, recordType string) error {
	rrset := new(RRset)
	rrset.Name = &name
	rrset.Type = &recordType
	rrset.ChangeType = String("DELETE")

	return r.patchRRset(domain, *rrset)
}

func (r *RecordsService) patchRRset(domain string, rrset RRset) error {
	rrset.Name = String(makeDomainCanonical(*rrset.Name))

	if *rrset.Type == "CNAME" {
		for i := range rrset.Records {
			rrset.Records[i].Content = String(makeDomainCanonical(*rrset.Records[i].Content))
		}
	}

	payload := RRsets{}
	payload.Sets = append(payload.Sets, rrset)

	req, err := r.client.newRequest("PATCH", fmt.Sprintf("servers/%s/zones/%s", r.client.VHost, trimDomain(domain)), nil, payload)
	if err != nil {
		return err
	}

	_, err = r.client.do(req, nil)
	return err
}
