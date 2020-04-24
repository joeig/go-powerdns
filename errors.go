package powerdns

import "fmt"

// Error structure with JSON API metadata
type Error struct {
	Status     string
	StatusCode int
	Message    string `json:"error"`
}

func (e Error) Error() string {
	return fmt.Sprintf("%v", e.Message)
}
