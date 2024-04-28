package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/spf13/pflag"
	"github.com/lorentzforces/ensure-path/internal/path_tools"
)

func main() {
	var anyPosition bool
	pflag.BoolVar(
		&anyPosition,
		"any-position",
		false,
		"Only verify that the item is present in the path, not necessarily first",
	)

	var useStdIn bool
	pflag.BoolVar(
		&useStdIn,
		"stdin",
		false,
		"Use standard input as the input, instead of the $PATH environment variable",
	)

	var helpRequested bool
	pflag.BoolVarP(
		&helpRequested,
		"help",
		"h",
		false,
		"Print this help message",
	)

	var keepEmpty bool
	pflag.BoolVarP(
		&keepEmpty,
		"keep-empty",
		"e",
		false,
		"Keep empty items in the input",
	)

	var deleteEntry string
	pflag.StringVarP(
		&deleteEntry,
		"delete",
		"d",
		"",
		"Delete items containing the provided string",
	)

	pflag.Parse()

	if helpRequested {
		printUsage()
		os.Exit(1)
	}

	args := pflag.Args()

	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Expected 1 string arg, but was given %d\n\n", len(args))
		printUsage()
		os.Exit(1)
	}

	entry := args[0]

	var pathString string
	var err error
	if useStdIn {
		pathString, err = getPathStdIn()
		if err != nil {
			failOut(err.Error())
		}
	} else {
		pathString = os.Getenv("PATH")
	}

	params := path_tools.EnsureParams{
		IncomingEntry: entry,
		Path: pathString,
		EnsureFirst: !anyPosition,
		RemoveEmpty: !keepEmpty,
		RemoveMatches: len(deleteEntry) > 0,
		MatchSeq: deleteEntry,
	}

	output := path_tools.EnsurePath(params)
	fmt.Println(output)
}

const maxTotalKilobytes = 10
const maxTotalBytes = 1024 * maxTotalKilobytes

// Pulls in standard input as a string until EOF is encountered. If there is no data in standard
// input, it will spin indefinitely (turns out piped input is kind of hard). We at least check to
// make sure STDIN is a pipe and not accidentally the terminal.
func getPathStdIn() (string, error) {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		err := errors.New("STDIN was not a pipe, which is probably not intended.")
		return "", err
	}

	stdIn := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	var out strings.Builder
	totalBytes := uint(0)

	for {
		n, err := stdIn.Read(buf[:cap(buf)])
		buf = buf[:n]
		totalBytes += uint(n)

		if totalBytes > maxTotalBytes {
			err := fmt.Errorf(
				"Standard input contained more than %d kB. This is probably not intended.",
				maxTotalKilobytes,
			)
			return out.String(), err
		}

		if n == 0 {
			if err == nil {
				continue
			}
			if errors.Is(err, io.EOF) {
				break
			}
		}

		if err != nil && !errors.Is(err, io.EOF) {
			return out.String(), nil
		}

		out.Write(buf)
	}

	return out.String(), nil
}

func printUsage() {
	fmt.Fprint(
		os.Stderr,
		`Usage of ensure-path:  ensure-path [OPTION]... NEWENTRY
Reads from the $PATH environment variable, returning that PATH value with NEWENTRY appearing once
at the beginning of its entries.

Options:
`,
	)
	pflag.PrintDefaults()
}

func failOut(msg string) {
	fmt.Fprintln(os.Stderr, "ERROR: " + msg)
	os.Exit(1)
}
