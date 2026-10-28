package domains

type Put[T any] struct {
	Value T
	Set   bool
}
