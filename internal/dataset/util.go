package dataset

func Filter[T any](slice []T, include func(T) bool) []T {
	var result []T

	for _, element := range slice {
		if include(element) {
			result = append(result, element)
		}
	}

	return result
}

func ToSliceOfAny[T any](in []T) []any {
	out := make([]any, len(in))

	for index, value := range in {
		out[index] = value
	}

	return out
}
