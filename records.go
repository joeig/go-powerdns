package powerdns

import "strings"

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

	if err != nil {
		return err
	}

	switch code := resp.StatusCode; {
	case code >= 400:
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return myError
	case code >= 200 && code <= 299:
		return nil
	default:
		return err
	}
}
