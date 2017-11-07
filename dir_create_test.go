// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func ExampleDirCreate() {
	root, cleanup, err := InitTempDir()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	err = os.Chdir(root)
	if err != nil {
		log.Fatal(err)
	}

	wd, _ := os.Getwd()
	log.Printf("Temp folder: %v\n", wd)

	dirs := `
2999-01-01T01:01:01Z 0777 dir_create_example

2001-01-01T01:01:01Z 0750 dir_create_example/a
2001-01-01T01:01:01Z 0750 dir_create_example/b

2002-01-01T01:01:01Z 0700 "dir_create_example/has\ttab"
2002-01-01T01:01:01Z 0700 "dir_create_example/\u263asmiles\u263a"
`

	reader := strings.NewReader(dirs)
	err = DirCreate(reader)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	files, err := ioutil.ReadDir("dir_create_example")
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	fmt.Printf("%v | %v | %v\n", files[1].ModTime().UTC(), files[1].Mode().Perm(), files[1].Name())
	// Output: 2001-01-01 01:01:01 +0000 UTC | -rwxr-x--- | b
}

func TestDirCreate(t *testing.T) {

	dirs := `
2001-01-01T01:01:01Z 0777 aaa
2009-01-01T01:01:01Z 0777 aaa/bbb

2002-01-01T01:01:01Z 0777 "has\ttab"
2002-01-01T01:01:01Z 0777 "\u10077heavy quoted\u10078"`

	expect := []struct {
		t time.Time
		p os.FileMode
		n string
	}{
		{time.Date(2001, time.January, 1, 1, 1, 1, 0, time.UTC), 0777, "aaa"},
		{time.Date(2002, time.January, 1, 1, 1, 1, 0, time.UTC), 0777, "has\ttab"},
		{time.Date(2002, time.January, 1, 1, 1, 1, 0, time.UTC), 0777, "\u10077heavy quoted\u10078"},
		{time.Date(2009, time.January, 1, 1, 1, 1, 0, time.UTC), 0777, "aaa/bbb"},
	}

	root, cleanup, err := InitTempDir()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer func(dir string) {
		os.Chdir(dir)
	}(wd)

	err = os.Chdir(root)
	if err != nil {
		t.Fatal(err)
	}

	reader := strings.NewReader(dirs)
	err = DirCreate(reader)
	if err != nil {
		t.Fatal(err)
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		t.Fatal(err)
		return
	}

	for i, fi := range files {
		if fi.Name() != expect[i].n {
			t.Errorf("Names mismatch, expected \"%v\", got \"%v\"", expect[i].n, fi.Name())
		}
		if fi.ModTime().UTC() != expect[i].t {
			t.Errorf("Times mismatch, expected \"%v\", got \"%v\" for \"%v\"", expect[i].t, fi.ModTime().UTC(), expect[i].n)
		}
		if fi.Mode().Perm() != expect[i].p {
			t.Errorf("Permissions mismatch, expected \"%v\", got \"%v\" for \"%v\"", expect[i].p, fi.Mode().Perm(), expect[i].n)
		}
	}

	f := expect[3].n
	fi, err := os.Stat(f)
	if err != nil {
		t.Fatal(err)
	}

	if fi.ModTime().UTC() != expect[3].t {
		t.Errorf("Times mismatch, expected \"%v\", got \"%v\" for \"%v\"", expect[3].t, fi.ModTime().UTC(), expect[3].n)
	}

	if fi.Mode().Perm() != expect[3].p {
		t.Errorf("Permissions mismatch, expected \"%v\", got \"%v\" for \"%v\"", expect[3].p, fi.Mode().Perm(), expect[3].n)
	}
}
