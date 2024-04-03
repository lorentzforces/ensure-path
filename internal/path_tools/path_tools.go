package path_tools

import (
	"fmt"
	"regexp"
	"strings"
)

const delimiter string = ":"

func EnsureOnce(entry, path string, removeEmpty bool) string {
	if len(path) == 0 && removeEmpty {
		return entry
	}
	if len(path) == 0 && !removeEmpty {
		return entry + ":"
	}

	pathEntries := checkPath(entry, path)

	if pathEntries.isEmpty() {
		return entry
	}
	if pathEntries.entryPresent() {
		return path
	}

	if removeEmpty {
		pathEntries.removeEmptyEntries()
	}

	return strings.Join(pathEntries.entries, delimiter)
}

func EnsureFirst(entry, path string, removeEmpty bool) string {
	if len(path) == 0 && removeEmpty {
		return entry
	}
	if len(path) == 0 && !removeEmpty {
		return entry + ":"
	}

	pathEntries := checkPath(entry, path)

	if pathEntries.isEmpty() {
		return entry
	}
	if pathEntries.entryFirst() {
		return path
	}

	if removeEmpty {
		pathEntries.removeEmptyEntries()
	}

	return strings.Join(pathEntries.entries, delimiter)
}

type splitPath struct {
	entries []string
	entryIndex int
}

const stringNotPresent int = -1

// Parses the provided path, splitting it into its component entries. It will record if and where
// the incoming entry was seen. Regardless, the returned struct will be modified to add the
// incoming entry at the beginning, and the consumer can decide whether to use this new set of
// entries or to discard it.
func checkPath(entry, path string) splitPath {
	rawEntries := strings.Split(path, delimiter)

	modifiedEntries := make([]string, 0, len(rawEntries) + 1)
	entryIndex := int(-1)
	modifiedEntries = append(modifiedEntries, entry)
	for i, existing := range rawEntries {
		if existing == entry && entryIndex == stringNotPresent {
			entryIndex = i
		} else {
			modifiedEntries = append(modifiedEntries, existing)
		}
	}

	return splitPath{
		entries: modifiedEntries,
		entryIndex: entryIndex,
	}
}

func (sp splitPath) isEmpty() bool {
	return len(sp.entries) == 1 && sp.entries[0] == ""
}

func (sp splitPath) entryPresent() bool {
	return sp.entryIndex > -1
}

func (sp splitPath) entryFirst() bool {
	return sp.entryIndex == 0
}

func (sp *splitPath) removeEmptyEntries() {
	newEntries := make([]string, 0, len(sp.entries))
	for _, existing := range sp.entries {
		whitespaceOnly, _ := regexp.MatchString("\\s", existing)
		if whitespaceOnly || len(existing) == 0 {
			continue
		}
		newEntries = append(newEntries, existing)
	}
	sp.entries = newEntries
}

func (sp splitPath) printDebug() string {
	var debugStr strings.Builder;
	debugStr.WriteString("[\n")
	for _, entry := range sp.entries {
		debugStr.WriteString(fmt.Sprintf("  \"%s\",\n", entry))
	}
	debugStr.WriteString("]\n")
	return debugStr.String()
}
