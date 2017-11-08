// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"
import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

type emptyErr struct {
	error
}

// TreeCreate parses a suplied Reader for the tree information
// and follows the instructions to create files and directories.
//
// The input has line records with three or four fields
// separated by one or more tabs. White space is trimmed on
// both ends of lines. Empty lines are skipped. The general
// line format is:
//
// <1. time>	<2. permissions>	<3. name> <4. optional content>
//
// Field 1: Time in RFC3339 format, as shown at
// https://golang.org/pkg/time/#RFC3339
//
// Field 2: Octal (required) representation of FileMode, as at
// https://golang.org/pkg/os/#FileMode
//
// Field 3: is the file or directory path to be created. If the
// first character of the path is a double-quote or a back-tick,
// then the path wil be passed through strconv.Unquote() function.
// It allows for using tab-containing or otherwise weird names.
// The quote or back-tick should be balanced at the end of
// the field.
//
// If the path in Field 3 ends with a forward slash, then it is
// treated as a directory, otherwise - as a regular file.
//
// Field 4: is optional content to be written into the file. It
// follows the same quotation rules as paths in Field 3.
// Directory entries ignore Field 4 if present.
//
// It is up to the caller to deal with conflicting file and
// directory names in the input. TreeCreate processes the input
// line-by-e or weird directory names.
//
// It is up to the caller to deal with conflicting file and
// directory names in the input. TreeCreate processes the
// input line-by-line and will return with error at a first
// problem it runs into.
func TreeCreate(config io.Reader) error {

	dirs := make([]*dirEntry, 0)

	scanner := bufio.NewScanner(config)
	for scanner.Scan() {

		mt, perm, name, content, err := parse(scanner.Text())
		if err != nil {
			if _, ok := err.(*emptyErr); ok {
				continue
			}
			return err
		}

		if name[len(name)-1] == '/' {
			name = name[:len(name)-1]

			if err = os.Mkdir(name, 0700); err != nil {
				return err
			}

			dirs = append(dirs, &dirEntry{name, perm, mt})
			continue
		}

		f, err := os.Create(name)
		if err != nil {
			return err
		}

		if len(content) > 0 {
			_, err = f.WriteString(content)
			if err != nil {
				return err
			}
		}

		err = f.Close()
		if err != nil {
			return err
		}

		err = os.Chmod(name, perm)
		if err != nil {
			return err
		}

		err = os.Chtimes(name, mt, mt)
		if err != nil {
			return err
		}
	}

	err := scanner.Err()
	if err != nil {
		return err
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		err := os.Chmod(dirs[i].name, dirs[i].perm)
		if err != nil {
			return err
		}

		err = os.Chtimes(dirs[i].name, dirs[i].time, dirs[i].time)
		if err != nil {
			return err
		}
	}

	return nil
}

var (
	re    = regexp.MustCompile(`^\s*([-0-9T:Z]+)\t+(0[0-7]{0,4})\t+([^\t]+)(\t+([^\t]+))?\s*$`)
	empty = regexp.MustCompile(`^\s*$`)
)

func parse(line string) (time.Time, os.FileMode, string, string, error) {

	if empty.MatchString(line) {
		return time.Time{}, 0, "", "", &emptyErr{}
	}

	parts := re.FindStringSubmatch(line)

	mt, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return time.Time{}, 0, "", "", err
	}
	mt = mt.Round(0)

	perm64, err := strconv.ParseUint(parts[2], 8, 32)
	if err != nil {
		return time.Time{}, 0, "", "", err
	}

	perm := os.FileMode(perm64)

	var path string
	if parts[3][0] == '`' || parts[3][0] == '"' {
		path, err = strconv.Unquote(parts[3])
		if err != nil {
			return time.Time{}, 0, "", "", err
		}
	} else {
		path = parts[3]
	}

	var content string
	if len(parts[5]) > 0 {

		if parts[5][0] == '`' || parts[5][0] == '"' {
			content, err = strconv.Unquote(parts[5])
			if err != nil {
				return time.Time{}, 0, "", "", err
			}
		} else {
			content = parts[5]
		}
	}

	return mt, perm, path, content, nil
}
