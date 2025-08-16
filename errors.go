package powerdns

// Error structure with JSON API metadata
type Error struct {
	StatusCode int    `json:"-"`
	Status     string `json:"-"`
	Message    string `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}
