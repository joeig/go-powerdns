package powerdns

import "fmt"

// ServersService handles communication with the servers related methods of the Client API
type ServersService service

// Server structure with JSON API metadata
type Server struct {
	Type       *string `json:"type,omitempty"`
	ID         *string `json:"id,omitempty"`
	DaemonType *string `json:"daemon_type,omitempty"`
	Version    *string `json:"version,omitempty"`
	URL        *string `json:"url,omitempty"`
	ConfigURL  *string `json:"config_url,omitempty"`
	ZonesURL   *string `json:"zones_url,omitempty"`
}

// List retrieves a list of Servers
func (s *ServersService) List() ([]Server, error) {
	req, err := s.client.newRequest("GET", "servers", nil, nil)
	if err != nil {
		return nil, err
	}

	servers := make([]Server, 0)
	_, err = s.client.do(req, &servers)
	return servers, err
}

// Get returns a certain Server
func (s *ServersService) Get(vHost string) (*Server, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("servers/%s", vHost), nil, nil)
	if err != nil {
		return nil, err
	}

	server := &Server{}
	_, err = s.client.do(req, &server)
	return server, err
}
