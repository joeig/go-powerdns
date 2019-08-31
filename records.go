package powerdns

import (
	"strings"
)

// RRset structure with JSON API metadata
type RRset struct {
	Name       string    `json:"name,omitempty"`
	Type       string    `json:"type,omitempty"`
	TTL        int       `json:"ttl,omitempty"`
	ChangeType string    `json:"changetype,omitempty"`
	Records    []Record  `json:"records,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}

// Record structure with JSON API metadata
type Record struct {
	Content  string `json:"content,omitempty"`
	Disabled bool   `json:"disabled"`
	SetPTR   bool   `json:"set-ptr,omitempty"`
}

// Comment structure with JSON API metadata
type Comment struct {
	Content    string `json:"content,omitempty"`
	Account    string `json:"account,omitempty"`
	ModifiedAt int    `json:"modified_at,omitempty"`
}

// RRsets structure with JSON API metadata
type RRsets struct {
	Sets []RRset `json:"rrsets,omitempty"`
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
	rrset.Name = fixDomainSuffix(rrset.Name)

	payload := RRsets{}
	payload.Sets = append(payload.Sets, rrset)

	myError := new(Error)

	zonesSling := z.PowerDNSHandle.makeSling()
	resp, err := zonesSling.New().Patch(strings.TrimRight(z.URL, ".")).BodyJSON(payload).Receive(nil, myError)
	if err != nil {
		return err
	}

	if err := handleAPIClientError(resp, &err, myError); err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return myError
	}

	return nil
}
