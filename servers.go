package powerdns

import (
	"strings"
)

type Server struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	DaemonType string `json:"daemon_type"`
	Version    string `json:"version"`
	URL        string `json:"url"`
	ConfigURL  string `json:"config_url"`
	ZonesURL   string `json:"zones_url"`
}

func (p *PowerDNS) GetServers() ([]Server, error) {
	serversSling := p.makeSling("servers")

	servers := make([]Server, 0)
	error := new(Error)
	resp, err := serversSling.New().Get("").Receive(&servers, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	return servers, err
}

func (p *PowerDNS) GetServer(serverID string) (Server, error) {
	serversSling := p.makeSling("servers/")

	server := Server{}
	error := new(Error)
	resp, err := serversSling.New().Get(serverID).Receive(&server, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return Server{}, error
	}

	return server, err
}
