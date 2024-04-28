package path_tools

import (
	"regexp"
	"strings"
)

const delimiter string = ":"

type EnsureParams struct {
	IncomingEntry string
	Path string
	EnsureFirst bool
	RemoveEmpty bool
	RemoveMatches bool
	MatchSeq string
}

func EnsurePath(params EnsureParams) string {
	entries := strings.Split(params.Path, delimiter)

	filters := make([]filterFunc, 0, 3)
	if params.RemoveEmpty {
		filters = append(filters, filterEmpty)
	}
	if params.RemoveMatches {
		filters = append(filters, filterBySubstring(params.MatchSeq))
	}
	if params.EnsureFirst {
		filters = append(filters, filterByString(params.IncomingEntry))
	}
	entries = filterEntries(entries, filters)

	entryFound := false
	for _, entry := range entries {
		if entry == params.IncomingEntry {
			entryFound = true
		}
	}

	if !entryFound {
		freshSlice := make([]string, 0, len(entries) + 1)
		freshSlice = append(freshSlice, params.IncomingEntry)
		freshSlice = append(freshSlice, entries...)
		entries = freshSlice
	}

	return strings.Join(entries, delimiter)
}

func filterEntries(entries []string, filters []filterFunc) []string {
	if len(filters) == 0 {
		return entries
	}

	results := make([]string, 0, len(entries))

	for _, entry := range entries {
		kept := true
		for _, filter := range filters {
			if !filter(entry) {
				kept = false
			}
		}

		if kept {
			results = append(results, entry)
		}
	}


	return results
}

// return true to keep entry in input list, false to remove
type filterFunc func(entry string) bool

func filterEmpty(entry string) bool {
	isWhitespace, _ := regexp.MatchString(`^\s*$`, entry)

	return !isWhitespace
}

func filterBySubstring(subStr string) filterFunc {
	return func(entry string) bool {
		return !strings.Contains(entry, subStr)
	}
}

func filterByString(str string) filterFunc {
	return func(entry string) bool {
		return str != entry
	}
}
