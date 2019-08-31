package powerdns

// Statistic structure with JSON API metadata
type Statistic struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// GetStatistics retrieves a list of Statistics
func (p *PowerDNS) GetStatistics() ([]Statistic, error) {
	statistics := make([]Statistic, 0)
	myError := new(Error)

	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost+"/statistics").Receive(&statistics, myError)

	if err := handleAPIClientError(resp, &err, myError); err != nil {
		return nil, err
	}

	return statistics, err
}
