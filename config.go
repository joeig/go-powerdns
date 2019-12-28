package powerdns

import "fmt"

// ConfigService handles communication with the zones related methods of the Client API
type ConfigService service

// ConfigSetting structure with JSON API metadata
type ConfigSetting struct {
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
	Value *string `json:"value,omitempty"`
}

// List retrieves a list of ConfigSettings
func (c *ConfigService) List() ([]ConfigSetting, error) {
	req, err := c.client.newRequest("GET", fmt.Sprintf("servers/%s/config", c.client.VHost), nil, nil)
	if err != nil {
		return nil, err
	}

	config := make([]ConfigSetting, 0)
	_, err = c.client.do(req, &config)
	return config, err
}
