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

// DirCreate parses a suplied Reader for directory information
// and follows the instructions to create directories.
//
// The file has line records with three space-separated fields.
//
// Field 1: Time in RFC3339 format, as shown at
// https://golang.org/pkg/time/#RFC3339
//
// Field 2: Octal (required) representation of FileMode, as at
// https://golang.org/pkg/os/#FileMode
//
// Field 3: a path to the directory to be create, until the end
// of the line. All characters are significant. If the first
// character of the path is a double-quote or a back-tick, then
// the path wil be passed through strconv.Unquote() function.
// It allows for creating Unicode or weird directory names.
//
// There should be no space before the first field
func DirCreate(config io.Reader) error {

	dirs := make([]*dirEntry, 0)

	scanner := bufio.NewScanner(config)
	for scanner.Scan() {

		mt, perm, name, err := parse(scanner.Text())
		if err != nil {
			if _, ok := err.(*emptyErr); ok {
				continue
			}
			return err
		}

		if err = os.Mkdir(name, 0700); err != nil {
			return err
		}

		dirs = append(dirs, &dirEntry{name, perm, mt})
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
	re    = regexp.MustCompile(`^([-0-9T:Z]+)\s+(0[0-7]{0,4})\s+(.*)$`)
	empty = regexp.MustCompile(`^\s*$`)
)

func parse(line string) (time.Time, os.FileMode, string, error) {

	if empty.MatchString(line) {
		return time.Time{}, 0, "", &emptyErr{}
	}

	parts := re.FindStringSubmatch(line)

	mt, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return time.Time{}, 0, "", err
	}

	perm64, err := strconv.ParseUint(parts[2], 8, 32)
	if err != nil {
		return time.Time{}, 0, "", err
	}

	perm := os.FileMode(perm64)

	var path string
	if parts[3][0] == '`' || parts[3][0] == '"' {
		path, err = strconv.Unquote(parts[3])
		if err != nil {
			return time.Time{}, 0, "", err
		}
	} else {
		path = parts[3]
	}

	return mt, perm, path, nil
}
