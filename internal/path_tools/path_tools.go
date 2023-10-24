package path_tools

import "strings"

const delimiter string = ":"

func EnsureOnce(entry, path string) string {
	pathEntries := checkPath(entry, path)

	if pathEntries.isEmpty() {
		return entry
	}
	if pathEntries.entryPresent() {
		return path
	}

	return strings.Join(pathEntries.entries, delimiter)
}

func EnsureFirst(entry, path string) string {
	pathEntries := checkPath(entry, path)

	if pathEntries.isEmpty() {
		return entry
	}
	if pathEntries.entryFirst() {
		return path
	}

	return strings.Join(pathEntries.entries, delimiter)
}

const stringNotPresent int = -1

type splitPath struct {
	entries []string
	entryIndex int
}

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
		if existing == entry && entryIndex == -1 {
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
