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

func getUniqueStrings(stringsList []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueStrings []string

	// Iterate over the strings list
	for _, str := range stringsList {
		if _, exists := uniqueMap[str]; !exists {
			uniqueMap[str] = true
			uniqueStrings = append(uniqueStrings, str)
		}
	}

	return uniqueStrings
}
