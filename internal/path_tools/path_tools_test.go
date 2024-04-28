package path_tools

import (
	"testing"
)

func TestEnsureOnceIsNoOpWhenEntryPresent(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "path/someDir:path/test:anotherPath/test",
	}

	result := EnsurePath(params)

	if result != params.Path {
		t.Errorf(
			"Path was updated when entry was already present.\nPath: %s\nOutput: %s",
			params.Path,
			result,
		)
	}
}

func TestEnsureOnceAddsEntry(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "path/someDir:anotherPath/test",
	}

	result := EnsurePath(params)

	expectedResult := params.IncomingEntry + ":" + params.Path
	if expectedResult != result {
		t.Errorf(
			"Result was expected to have new entry \"%s\" at the beginning.\nResult: %s",
			params.IncomingEntry,
			result,
		)
	}
}

func TestEnsureFirstIsNoOpWhenEntryFirst(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "path/test:path/someDir:anotherPath/test",
		EnsureFirst: true,
	}

	result := EnsurePath(params)

	if result != params.Path {
		t.Errorf(
			"Path was updated when entry was already present at beginning.\nPath: %s\nOutput: %s",
			params.Path,
			result,
		)
	}
}

func TestEnsureFirstAddsEntryWhenAbsent(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "path/someDir:anotherPath/test",
		EnsureFirst: true,
	}

	result := EnsurePath(params)

	expectedResult := params.IncomingEntry + ":" + params.Path
	if expectedResult != result {
		t.Errorf(
			"Result was expected to have new entry \"%s\" at the beginning.\nResult: %s",
			params.IncomingEntry,
			result,
		)
	}
}

func TestEnsureFirstMovesEntryWhenNotFirst(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "path/someDir:path/test:anotherPath/test",
		EnsureFirst: true,
	}

	result := EnsurePath(params)

	expectedResult := "path/test:path/someDir:anotherPath/test"
	if expectedResult != result {
		t.Errorf(
			"Result was expected to have entry \"%s\" at the beginning.\nResult: %s",
			params.IncomingEntry,
			result,
		)
	}
}

func TestEmptyPathWithoutEmptyRemovalAddsDelimiter(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "",
	}

	result := EnsurePath(params)

	expectedResult := params.IncomingEntry + ":"
	if expectedResult != result {
		t.Errorf(
			"Expected result to retain implicit empty entry, should look like \"%s\".\nResult: %s",
			expectedResult,
			result,
		)
	}
}

func TestEmptyPathWithEmptyRemovalIsLiteralEntry(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "",
		RemoveEmpty: true,
	}

	result := EnsurePath(params)

	if result != params.IncomingEntry {
		t.Errorf(
			"Expected result to be identical to input (no delimiter).\nResult: %s",
			result,
		)
	}
}

func TestEmptyEntriesAreRemoved(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "path/someDir::path/test",
		RemoveEmpty: true,
		EnsureFirst: true,
	}

	result := EnsurePath(params)

	expectedResult := "path/test:path/someDir"
	if result != expectedResult {
		t.Errorf(
			"Expected result to no longer have an empty entry \"::\".\nResult: %s",
			result,
		)
	}
}

func TestEmptyEntriesAreKept(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/test",
		Path: "path/someDir::path/test",
		EnsureFirst: true,
	}

	result := EnsurePath(params)

	expectedResult := "path/test:path/someDir:"
	if result != expectedResult {
		t.Errorf(
			"Expected result to have an empty entry \"::\".\nResult: %s",
			result,
		)
	}
}

func TestMatchedPathsAreDeleted(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/anotherOne",
		Path: "path/test:path/deleteme:path/someDir",
		RemoveMatches: true,
		MatchSeq: "deleteme",
	}

	result := EnsurePath(params)

	expectedResult := "path/anotherOne:path/test:path/someDir"
	if result != expectedResult {
		t.Errorf(
			"Expected result to not have \"deleteme\" entry.\nResult: %s",
			result,
		)
	}
}

func TestCanAddEntryThatMatchesDeleted(t *testing.T) {
	params := EnsureParams{
		IncomingEntry: "path/deleteme",
		Path: "path/test:path/deleteme:path/someDir",
		RemoveMatches: true,
		MatchSeq: "deleteme",
	}

	result := EnsurePath(params)

	expectedResult := "path/deleteme:path/test:path/someDir"
	if result != expectedResult {
		t.Errorf(
			"Expected result to have single \"path/deleteme\" entry at beginning.\nResult: %s",
			result,
		)
	}
}
