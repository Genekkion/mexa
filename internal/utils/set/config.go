package set

// Opt is the type for setting a config option for the set to be created.
type Opt[T comparable] func(*Set[T])

// WithSlice Adds elements from an initial slice of elements.
func WithSlice[T comparable](slice []T) Opt[T] {
	return func(s *Set[T]) {
		for _, v := range slice {
			s.m[v] = emptyValue
		}
	}
}
