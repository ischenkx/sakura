package util

func MakePtr[T any](data T) *T {
	return &data
}
