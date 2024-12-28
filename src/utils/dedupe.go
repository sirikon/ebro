package utils

func Dedupe[T comparable](slice []T) []T {
	result := []T{}
	itemExists := map[T]bool{}
	for _, item := range slice {
		if _, ok := itemExists[item]; !ok {
			result = append(result, item)
			itemExists[item] = true
		}
	}
	return result
}
