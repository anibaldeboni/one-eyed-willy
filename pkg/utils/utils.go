package utils

func contains[T comparable](slice []T, element T) bool {
	for _, a := range slice {
		if a == element {
			return true
		}
	}
	return false
}

// # MAY I HAVE YOUR ATTENTION, PLEASE! #
//
// DO NOT USE FOR LARGE SLICES
// Code is not optimized for performance
func IsSubSlice[T comparable](slice []T, subslice []T) bool {
	if len(slice) < len(subslice) {
		return false
	}
	for _, e := range slice {
		if !contains(subslice, e) {
			return false
		}
	}
	return true
}
