package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	if err := run(); err == flag.ErrHelp {
		usage()
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	// Parse arguments.
	fs := flag.NewFlagSet("fslice", flag.ContinueOnError)
	fs.Usage = func() {}
	out := fs.String("o", "", "output file")
	start := fs.String("start", "", "starting delimiter")
	end := fs.String("end", "", "ending delimiter")
	header := fs.String("header", "", "header line")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}

	// Validate arguments.
	paths := fs.Args()
	if len(paths) == 0 {
		return errors.New("path required")
	} else if *start == "" {
		return errors.New("starting delimiter required")
	} else if *end == "" {
		return errors.New("ending delimiter required")
	}

	// Clean up arguments.
	*start = strings.TrimSpace(*start)
	*end = strings.TrimSpace(*end)

	// Process each path.
	var buf bytes.Buffer
	for _, path := range paths {
		if err := process(&buf, path, *start, *end, *header); err != nil {
			return err
		}
	}

	// Write to STDOUT if no file specified.
	if *out == "" {
		buf.WriteTo(os.Stdout)
		return nil
	}

	// If writing to a file, ensure it hasn't changed.
	data := buf.Bytes()
	if b, err := ioutil.ReadFile(*out); err != nil && !os.IsNotExist(err) {
		return err
	} else if bytes.Equal(data, b) {
		return nil // unchanged, exit
	}

	return ioutil.WriteFile(*out, data, 0666)
}

func process(w io.Writer, path, start, end, header string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var inBlock bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if we change state of in or out of block.
		if inBlock && strings.TrimSpace(line) == end {
			inBlock = false
			if header != "" {
				fmt.Fprintln(w, "")
			}
			continue
		} else if !inBlock && strings.TrimSpace(line) == start {
			inBlock = true
			if header != "" {
				fmt.Fprintln(w, strings.Replace(header, "$FILENAME", path, -1))
			}
			continue
		}

		// Print all lines while we're in a start/end block.
		if inBlock {
			fmt.Fprintln(w, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func usage() {
	fmt.Print(strings.TrimSpace(`
fslice is a small utility for extracting delimited sections of a file.

Usage:

	fslice [arguments] PATH [PATH]

Arguments:
	-start DELIM
	    Delimiting line that begins an output block.

	-end DELIM
	    Delimiting line that ends an output block.

	-header STR
	    String printed at the top of each output block. The special
	    $FILENAME variable can be used to print out the filename.
`) + "\n\n")
}
