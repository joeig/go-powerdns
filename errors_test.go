package powerdns

import "testing"

func TestError(t *testing.T) {
	myError := &Error{Message: "foo"}
	if myError.Error() != "foo" {
		t.Error("Error method returns invalid format")
	}
}
