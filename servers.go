package powerdns

// Server structure with JSON API metadata
type Server struct {
	Type       string `json:"type,omitempty"`
	ID         string `json:"id,omitempty"`
	DaemonType string `json:"daemon_type,omitempty"`
	Version    string `json:"version,omitempty"`
	URL        string `json:"url,omitempty"`
	ConfigURL  string `json:"config_url,omitempty"`
	ZonesURL   string `json:"zones_url,omitempty"`
}

// GetServers retrieves a list of Servers
func (p *PowerDNS) GetServers() ([]Server, error) {
	servers := make([]Server, 0)
	myError := new(Error)

	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers").Receive(&servers, myError)

	if err := handleAPIClientError(resp, &err, myError); err != nil {
		return nil, err
	}

	return servers, err
}

// GetServer returns a certain Server
func (p *PowerDNS) GetServer() (*Server, error) {
	server := &Server{}
	myError := new(Error)

	serversSling := p.makeSling()
	resp, err := serversSling.New().Get("servers/"+p.VHost).Receive(&server, myError)

	if err := handleAPIClientError(resp, &err, myError); err != nil {
		return nil, err
	}

	return server, err
}
