package powerdns

import "strings"

type Statistic struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (p *PowerDNS) GetStatistics(serverID string) ([]Statistic, error) {
	serversSling := p.makeSling("servers/")

	statistics := make([]Statistic, 0)
	error := new(Error)
	resp, err := serversSling.New().Get(serverID+"/statistics").Receive(&statistics, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	return statistics, err
}
