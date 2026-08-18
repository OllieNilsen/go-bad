package main

// Reverse returns the elements of slice in reverse order.
// Do not use any standard library functions to achieve this.
func Reverse(slice []int) []int {
	// TODO: implement
	return nil
}

// Intersect returns the elements present in both a and b.
func Intersect(a, b []int) []int {
	result := []int{}
	for _, x := range a {
		for _, y := range b {
			if x == y {
				result = append(result, x)
			}
		}
	}
	return result
}
