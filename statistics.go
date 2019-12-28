package powerdns

import (
	"fmt"
	"net/url"
)

// StatisticsService handles communication with the statistics related methods of the Client API
type StatisticsService service

// Statistic structure with JSON API metadata
type Statistic struct {
	Name *string `json:"name,omitempty"`
	Type *string `json:"type,omitempty"`

	// Contrary to the authoritative API specification, the "size" field has actually been implemented as string instead of integer.
	Size *string `json:"size,omitempty"`

	// The "value" field contains either a string or a list of objects, depending on the "type".
	Value interface{} `json:"value,omitempty"`
}

// List retrieves a list of Statistics
func (s *StatisticsService) List() ([]Statistic, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("servers/%s/statistics", s.client.VHost), nil, nil)
	if err != nil {
		return nil, err
	}

	statistics := make([]Statistic, 0)
	_, err = s.client.do(req, &statistics)

	return statistics, err
}

// Get retrieves certain Statistics
func (s *StatisticsService) Get(statisticName string) ([]Statistic, error) {
	query := url.Values{}
	query.Add("statistic", statisticName)
	req, err := s.client.newRequest("GET", fmt.Sprintf("servers/%s/statistics", s.client.VHost), &query, nil)
	if err != nil {
		return nil, err
	}

	statistics := make([]Statistic, 0)
	_, err = s.client.do(req, &statistics)

	return statistics, err
}
