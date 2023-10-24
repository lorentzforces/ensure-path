package path_tools

import "testing"

func TestEnsureOnceIsNoOpWhenEntryPresent(t *testing.T) {
	entry := "path/test"
	path := "path/someDir:path/test:anotherPath/test"

	result := EnsureOnce(entry, path)

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

	result := EnsureOnce(entry, path)

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

	result := EnsureFirst(entry, path)

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

	result := EnsureFirst(entry, path)

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

	result := EnsureFirst(entry, path)

	expectedResult := "path/test:path/someDir:anotherPath/test"
	if expectedResult != result {
		t.Errorf(
			"Result was expected to have entry \"%s\" at the beginning.\nResult: %s",
			entry,
			result,
		)
	}
}
