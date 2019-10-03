package powerdns

import "fmt"

// StatisticsService handles communication with the statistics related methods of the Client API
type StatisticsService service

// Statistic structure with JSON API metadata
type Statistic struct {
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
	Value *string `json:"value,omitempty"`
}

// List retrieves a list of Statistics
func (s *StatisticsService) List() ([]Statistic, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("servers/%s/statistics", s.client.VHost), nil)
	if err != nil {
		return nil, err
	}

	statistics := make([]Statistic, 0)
	_, err = s.client.do(req, &statistics)

	return statistics, err
}
