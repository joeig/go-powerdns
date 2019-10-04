package powerdns

import "testing"

func TestBool(t *testing.T) {
	source := true
	if *Bool(source) != source {
		t.Error("Invalid return value")
	}
}

func TestBoolValue(t *testing.T) {
	source := true
	if BoolValue(&source) != source {
		t.Error("Invalid return value")
	}
	if BoolValue(nil) != false {
		t.Error("Unexpected return value")
	}
}

func TestUint32(t *testing.T) {
	source := uint32(1337)
	if *Uint32(source) != source {
		t.Error("Invalid return value")
	}
}

func TestUint32Value(t *testing.T) {
	source := uint32(1337)
	if Uint32Value(&source) != source {
		t.Error("Invalid return value")
	}
	if Uint32Value(nil) != 0 {
		t.Error("Unexpected return value")
	}
}

func TestUint64(t *testing.T) {
	source := uint64(1337)
	if *Uint64(source) != source {
		t.Error("Invalid return value")
	}
}

func TestUint64Value(t *testing.T) {
	source := uint64(1337)
	if Uint64Value(&source) != source {
		t.Error("Invalid return value")
	}
	if Uint64Value(nil) != 0 {
		t.Error("Unexpected return value")
	}
}

func TestString(t *testing.T) {
	source := "foo"
	if *String(source) != source {
		t.Error("Invalid return value")
	}
}

func TestStringValue(t *testing.T) {
	source := "foo"
	if StringValue(&source) != source {
		t.Error("Invalid return value")
	}
	if StringValue(nil) != "" {
		t.Error("Unexpected return value")
	}
}
