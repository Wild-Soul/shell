package main

import (
	"strings"
)

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

func parseCommand(cmd string) (string, []string) {
	cmd = strings.Trim(cmd, " ") // remove starting and ending white spaces
	i := 0
	seenQuote := false // have we seen an opening quote before?
	args := []string{}
	var commandBuilder, argBuilder strings.Builder

	for {
		if i >= len(cmd) {
			break
		}
		if cmd[i] == ' ' {
			i += 1
			break
		}
		commandBuilder.WriteRune(rune(cmd[i]))
		i += 1
	}

	// now i points to starting of argument.
	for {
		if i >= len(cmd) {
			break
		}
		if cmd[i] == '\'' { // if it's a quote then toggle state.
			seenQuote = !seenQuote
		} else {
			if cmd[i] == ' ' {
				if seenQuote {
					argBuilder.WriteRune(rune(cmd[i]))
				} else {
					if argBuilder.Len() > 0 {
						args = append(args, argBuilder.String())
						argBuilder.Reset()
					}
				}
			} else {
				argBuilder.WriteRune(rune(cmd[i]))
			}
		}
		i += 1
	}

	if argBuilder.Len() > 0 {
		args = append(args, argBuilder.String())
		argBuilder.Reset()
	}

	return commandBuilder.String(), args
}
