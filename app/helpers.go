package main

func mapFunction[T any, U any](collection []T, f func(T) U) []U {
	var result []U

	for _, ele := range collection {
		result = append(result, f(ele))
	}

	return result
}

func filterFunction[T any](collection []T, f func(T) bool) []T {
	var result []T

	for _, ele := range collection {
		if f(ele) {
			result = append(result, ele)
		}
	}

	return result
}
