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

func (p *PowerDNS) GetServers() (*[]Server, error) {
	servers := make([]Server, 0)
	error := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers").Receive(&servers, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return nil, error
	}

	return &servers, err
}

func (p *PowerDNS) GetServer() (*Server, error) {
	server := &Server{}
	error := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost).Receive(&server, error)

	if err == nil && resp.StatusCode >= 400 {
		error.Message = strings.Join([]string{resp.Status, error.Message}, " ")
		return &Server{}, error
	}

	return server, err
}
