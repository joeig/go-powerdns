package powerdns

import (
	"strings"
)

// Server structure with JSON API metadata
type Server struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	DaemonType string `json:"daemon_type"`
	Version    string `json:"version"`
	URL        string `json:"url"`
	ConfigURL  string `json:"config_url"`
	ZonesURL   string `json:"zones_url"`
}

// GetServers retrieves a list of Servers
func (p *PowerDNS) GetServers() ([]Server, error) {
	servers := make([]Server, 0)
	myError := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers").Receive(&servers, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return nil, myError
	}

	return servers, err
}

// GetServer returns a certain Server
func (p *PowerDNS) GetServer() (*Server, error) {
	server := &Server{}
	myError := new(Error)
	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost).Receive(&server, myError)

	if err == nil && resp.StatusCode >= 400 {
		myError.Message = strings.Join([]string{resp.Status, myError.Message}, " ")
		return &Server{}, myError
	}

	return server, err
}
