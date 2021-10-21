package conv

// BoolPtrFrom converts a `bool` to `*bool`
func BoolPtrFrom(n bool) *bool {
	return &n
}

// BoolFromPtr converts a `*bool` to `bool`
func BoolFromPtr(n *bool) bool {
	if n == nil {
		return false
	}

	return *n
}
