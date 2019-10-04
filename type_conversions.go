package powerdns

// Bool is a helper function that allocates a new bool value to store v and returns a pointer to it.
func Bool(v bool) *bool {
	return &v
}

// BoolValue is a helper function that returns the value of a bool pointer or false.
func BoolValue(v *bool) bool {
	if v != nil {
		return *v
	}
	return false
}

// Uint32 is a helper function that allocates a new uint32 value to store v and returns a pointer to it.
func Uint32(v uint32) *uint32 {
	return &v
}

// Uint32Value is a helper function that returns the value of a bool pointer or 0.
func Uint32Value(v *uint32) uint32 {
	if v != nil {
		return *v
	}
	return 0
}

// Uint64 is a helper function that allocates a new uint64 value to store v and returns a pointer to it.
func Uint64(v uint64) *uint64 {
	return &v
}

// Uint64Value is a helper function that returns the value of a bool pointer or 0.
func Uint64Value(v *uint64) uint64 {
	if v != nil {
		return *v
	}
	return 0
}

// String is a helper function that allocates a new string value to store v and returns a pointer to it.
func String(v string) *string {
	return &v
}

// StringValue is a helper function that returns the value of a bool pointer or "".
func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
