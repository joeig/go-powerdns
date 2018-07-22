package powerdns

import (
	"strings"
)

type Statistic struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (p *PowerDNS) GetStatistics() ([]Statistic, error) {
	statistics := make([]Statistic, 0)
	error := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost+"/statistics").Receive(&statistics, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	return statistics, err
}
