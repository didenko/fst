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

	testRootDir, cleanup, err := InitTempDir()
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

	err = os.Chdir(testRootDir)
	if err != nil {
		t.Fatal(err)
	}

	dirs := `
2001-01-01T01:01:01Z 0050 aaa
2999-01-01T01:01:01Z 0700 aaa/bbb

2002-01-01T01:01:01Z 0700 "has\ttab"
2002-01-01T01:01:01Z 0700 "\u10077heavy quoted\u10078"`

	reader := strings.NewReader(dirs)
	err = DirCreate(reader)
	if err != nil {
		t.Fatal(err)
	}
}
