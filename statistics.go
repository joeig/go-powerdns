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
	myError := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost+"/statistics").Receive(&statistics, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return nil, myError
	}

	return statistics, err
}
