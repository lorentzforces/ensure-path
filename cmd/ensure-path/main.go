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
	var ensureFirst bool
	flag.BoolVar(
		&ensureFirst,
		"first",
		false,
		"Verify that the item is FIRST in the path, and add or move to beginning otherwise",
	)
	var useEnvVar bool
	flag.BoolVar(
		&useEnvVar,
		"from-env",
		false,
		"Use the value of the $PATH environment variable as the input",
	)
	flag.Parse()

	args := flag.Args()

	if len(args) != 1 {
		failOut(fmt.Sprintf("Expected 1 string arg, but was given %d", len(args)))
	}

	entry := args[0]

	var pathString string
	var err error
	if useEnvVar {
		pathString = os.Getenv("PATH")
	} else {
		pathString, err = getPathStdIn()
		if err != nil {
			failOut(err.Error())
		}
	}

	var output string
	if ensureFirst {
		output = path_tools.EnsureFirst(entry, pathString)
	} else {
		output = path_tools.EnsureOnce(entry, pathString)
	}


	fmt.Println(output)
}

const maxTotalKilobytes = 10
const maxTotalBytes = 1024 * maxTotalKilobytes

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
