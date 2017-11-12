// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

func ExampleTreeCreate() {

	dirMark := func(fi os.FileInfo) string {
		if fi.IsDir() {
			return "/"
		}
		return ""
	}

	_, cleanup, err := TempInitChdir()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	dirs := `
		2001-01-01T01:01:01Z	0750	a/
		2001-01-01T01:01:01Z	0750	b/
		2001-01-01T01:01:01Z	0700	c.txt	"This is a two line\nfile with\ta tab\n"
		2001-01-01T01:01:01Z	0700	d.txt	No need to quote a single line without tabs

		2002-01-01T01:01:01Z	0700	"has\ttab/"
		2002-01-01T01:01:01Z	0700	"has\ttab/e.mb"	"# Markdown...\n\n... also ***possible***\n"

		2002-01-01T01:01:01Z	0700	"\u263asmiles\u263a/"
	`

	reader := strings.NewReader(dirs)
	err = TreeCreate(reader)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

	fmt.Printf(
		"%v | %v | %s%s\n",
		files[1].ModTime().UTC(),
		files[1].Mode().Perm(),
		files[1].Name(),
		dirMark(files[1]),
	)

	fmt.Printf(
		"%v | %v | %s%s\n",
		files[2].ModTime().UTC(),
		files[2].Mode().Perm(),
		files[2].Name(),
		dirMark(files[2]),
	)

	// Output:
	// 2001-01-01 01:01:01 +0000 UTC | -rwxr-x--- | b/
	// 2001-01-01 01:01:01 +0000 UTC | -rwx------ | c.txt
}

type tcase struct {
	t time.Time
	p os.FileMode
	n string
	c string
}

func match(t *testing.T, f *tcase, fi os.FileInfo) {

	if !fi.ModTime().Equal(f.t) {
		t.Errorf("Times mismatch, expected \"%v\", got \"%v\" for \"%v\"", f.t, fi.ModTime().UTC(), f.n)
	}
	if fi.Mode().Perm() != f.p {
		t.Errorf("Permissions mismatch, expected \"%v\", got \"%v\" for \"%v\"", f.p, fi.Mode().Perm(), f.n)
	}

	if f.c == "" {
		return
	}

	byteContent, err := ioutil.ReadFile(f.n)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal([]byte(f.c), byteContent) {
		t.Errorf("Content mismatch, expected \"%v\", got \"%v\" for \"%v\"", f.c, byteContent, f.n)
	}
}

func TestTreeCreate(t *testing.T) {

	dirs := `
		2001-01-01T01:01:01Z	0150	aaa/
		2099-01-01T01:01:01Z	0700	aaa/bbb/

		2001-01-01T01:01:01Z	0700	c.txt	"This is a two line\nfile with\ta tab\n"
		2001-01-01T01:01:01Z	0700	d.txt	No need to quote a single line without tabs

		2002-01-01T01:01:01Z	0700	"has\ttab/"
		2002-01-01T01:01:01Z	0700	"has\ttab/e.mb"	"# Markdown...\n\n... also ***possible***\n"

		2002-01-01T01:01:01Z	0700	"\u10077heavy quoted\u10078/"`

	expect := []tcase{
		{time.Date(2001, time.January, 1, 1, 1, 1, 0, time.UTC), 0150, "aaa", ""},
		{time.Date(2001, time.January, 1, 1, 1, 1, 0, time.UTC), 0700, "c.txt", ""},
		{time.Date(2001, time.January, 1, 1, 1, 1, 0, time.UTC), 0700, "d.txt", "No need to quote a single line without tabs"},
		{time.Date(2002, time.January, 1, 1, 1, 1, 0, time.UTC), 0700, "has\ttab", ""},
		{time.Date(2002, time.January, 1, 1, 1, 1, 0, time.UTC), 0700, "\u10077heavy quoted\u10078", ""},
		{time.Date(2099, time.January, 1, 1, 1, 1, 0, time.UTC), 0700, "aaa/bbb", ""},
	}

	_, cleanup, err := TempInitChdir()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	reader := strings.NewReader(dirs)
	err = TreeCreate(reader)
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
		match(t, &expect[i], fi)
	}

	for _, tc := range expect {
		f := tc.n
		fi, err := os.Stat(f)
		if err != nil {
			t.Fatal(err)
		}

		match(t, &tc, fi)
	}
}
