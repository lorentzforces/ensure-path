package path_tools

import "testing"

func TestEnsureOnceIsNoOpWhenEntryPresent(t *testing.T) {
	entry := "path/test"
	path := "path/someDir:path/test:anotherPath/test"

	result := EnsureOnce(entry, path, false)

	if result != path {
		t.Errorf(
			"Path was updated when entry was already present.\nPath: %s\nOutput: %s",
			path,
			result,
		)
	}
}

func TestEnsureOnceAddsEntry(t *testing.T) {
	entry := "path/test"
	path := "path/someDir:anotherPath/test"

	result := EnsureOnce(entry, path, false)

	expectedResult := entry + ":" + path
	if expectedResult != result {
		t.Errorf(
			"Result was expected to have new entry \"%s\" at the beginning.\nResult: %s",
			entry,
			result,
		)
	}
}

func TestEnsureFirstIsNoOpWhenEntryFirst(t *testing.T) {
	entry := "path/test"
	path := "path/test:path/someDir:anotherPath/test"

	result := EnsureFirst(entry, path, false)

	if result != path {
		t.Errorf(
			"Path was updated when entry was already present at beginning.\nPath: %s\nOutput: %s",
			path,
			result,
		)
	}
}

func TestEnsureFirstAddsEntryWhenAbsent(t *testing.T) {
	entry := "path/test"
	path := "path/someDir:anotherPath/test"

	result := EnsureFirst(entry, path, false)

	expectedResult := entry + ":" + path
	if expectedResult != result {
		t.Errorf(
			"Result was expected to have new entry \"%s\" at the beginning.\nResult: %s",
			entry,
			result,
		)
	}
}

func TestEnsureFirstMovesEntryWhenNotFirst(t *testing.T) {
	entry := "path/test"
	path := "path/someDir:path/test:anotherPath/test"

	result := EnsureFirst(entry, path, false)

	expectedResult := "path/test:path/someDir:anotherPath/test"
	if expectedResult != result {
		t.Errorf(
			"Result was expected to have entry \"%s\" at the beginning.\nResult: %s",
			entry,
			result,
		)
	}
}

func TestEmptyPathWithoutEmptyRemovalAddsDelimiter(t *testing.T) {
	entry := "path/test"
	path := ""

	result := EnsureOnce(entry, path, false)

	expectedResult := entry + ":"
	if expectedResult != result {
		t.Errorf(
			"Expected result to retain implicit empty entry, should look like \"%s\".\nResult: %s",
			expectedResult,
			result,
		)
	}
}

func TestEmptyPathWithEmptyRemovalIsLiteralEntry(t *testing.T) {
	entry := "path/test"
	path := ""

	result := EnsureOnce(entry, path, true)

	if result != entry {
		t.Errorf(
			"Expected result to be identical to input (no delimiter).\nResult: %s",
			result,
		)
	}
}

func TestEmptyEntriesAreRemoved(t *testing.T) {
	entry := "path/test"
	path := "path/someDir::path/test"

	result := EnsureFirst(entry, path, true)

	expectedResult := "path/test:path/someDir"
	if result != expectedResult {
		t.Errorf(
			"Expected result to no longer have an empty entry \"::\".\nResult: %s",
			result,
		)
	}
}


func TestEmptyEntriesAreKept(t *testing.T) {
	entry := "path/test"
	path := "path/someDir::path/test"

	result := EnsureFirst(entry, path, false)

	expectedResult := "path/test:path/someDir:"
	if result != expectedResult {
		t.Errorf(
			"Expected result to have an empty entry \"::\".\nResult: %s",
			result,
		)
	}
}
