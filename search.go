package powerdns

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// SearchService handles communication with the search related methods of the Client API
type SearchService service

// SearchResult structure with JSON API metadata
type SearchResult struct {
	Content    *string `json:"content,omitempty"`
	Disabled   *bool   `json:"disabled,omitempty"`
	Name       *string `json:"name,omitempty"`
	ObjectType *string `json:"object_type,omitempty"`
	ZoneID     *string `json:"zone_id,omitempty"`
	Zone       *string `json:"zone,omitempty"`
	Type       *string `json:"type,omitempty"`
	TTL        *uint32 `json:"ttl,omitempty"`
}

// SearchObjectType string type
type SearchObjectType string

const (
	// SearchObjectTypeAll searches for all object types
	SearchObjectTypeAll SearchObjectType = "all"
	// SearchObjectTypeZone searches for zones
	SearchObjectTypeZone SearchObjectType = "zone"
	// SearchObjectTypeRecord searches for records
	SearchObjectTypeRecord SearchObjectType = "record"
	// SearchObjectTypeComment searches for comments
	SearchObjectTypeComment SearchObjectType = "comment"
)

// Search searches the PowerDNS server for data matching the query.
// The query supports wildcards: * for multiple characters, ? for a single character.
// The max parameter limits the number of returned results.
// The objectType parameter filters results by type (all, zone, record, comment).
func (s *SearchService) Search(ctx context.Context, query string, max int, objectType SearchObjectType) ([]SearchResult, error) {
	q := url.Values{}
	q.Add("q", query)
	q.Add("max", strconv.Itoa(max))
	if objectType != "" {
		q.Add("object_type", string(objectType))
	}

	req, err := s.client.newRequest(ctx, http.MethodGet, path.Join("servers", s.client.VHost, "search-data"), &q, nil)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0)
	_, err = s.client.do(req, &results)
	return results, err
}
