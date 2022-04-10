package util

func Map[T, R any](mapper func(T) R, slice []T) []R {
	if slice == nil {
		return nil
	}

	re := make([]R, len(slice))
	for i := range slice {
		re = append(re, mapper(slice[i]))
	}

	return re
}
