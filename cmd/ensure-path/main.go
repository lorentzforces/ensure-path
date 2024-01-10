package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"github.com/lorentzforces/ensure-path/internal/path_tools"
)

func main() {
	var anyPosition bool
	flag.BoolVar(
		&anyPosition,
		"any-position",
		false,
		"Only verify that the item is present in the path, not necessarily first.",
	)

	var useStdIn bool
	flag.BoolVar(
		&useStdIn,
		"stdin",
		false,
		"Use standard input as the input, instead of STDIN",
	)

	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		failOut(fmt.Sprintf("Expected 1 string arg, but was given %d", len(args)))
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

	var output string
	if anyPosition {
		output = path_tools.EnsureOnce(entry, pathString)
	} else {
		output = path_tools.EnsureFirst(entry, pathString)
	}

	fmt.Println(output)
}

const maxTotalKilobytes = 10
const maxTotalBytes = 1024 * maxTotalKilobytes

// Pulls piped standard input in. You can check
func getPathStdIn() (string, error) {
	stdIn := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	var out strings.Builder
	totalBytes := uint(0)

	for {
		n, err := stdIn.Read(buf[:cap(buf)])
		buf = buf[:n]
		totalBytes += uint(n)

		if totalBytes > maxTotalBytes {
			err := errors.New(
				"Standard input contained more than %d kB. This is probably not intended.",
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

func failOut(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
