// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func ExampleTempInitDir() {

	root, cleanup, err := TempInitDir()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestTempInitDir(t *testing.T) {

	// Get the values and create the test root dir to be tested
	testRootDir, cleanup, err := TempInitDir()
	if err != nil {
		t.Fatal(err)
	}

	// Check that the returned testRootDir is Stat-able
	testRootInfo, err := os.Stat(testRootDir)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the returned testRootDir is a directory
	if !testRootInfo.IsDir() {
		t.Errorf("Returned temporary path \"%s\" is not a directory", testRootDir)
	}

	// run the resulting cleaup function for tests below
	cleanup()

	// Check that the testRootDir does not exist after the cleanup
	_, err = os.Stat(testRootDir)
	_, ok := err.(*os.PathError)
	if err == nil && !ok {
		t.Fatalf("Temporary directory \"%s\" remained after cleanup", testRootDir)
	}
}

func TestTempInitChdir(t *testing.T) {

	// Capture the old workdir
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Get the values and create the test root dir to be tested
	old, cleanup, err := TempInitChdir()
	if err != nil {
		t.Fatal(err)
	}

	if origWD != old {
		t.Fatalf("Got \"%s\" as an old directory instead of the expected \"%s\"\n", old, origWD)
	}

	// Capture the temporary workdir
	tempWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Check that we are in an empty directory
	files, err := ioutil.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) > 0 {
		t.Fatalf("The current, supposedly new, directory \"%s\" is not empty\n", tempWD)
	}

	// run the resulting cleaup function for tests below
	cleanup()

	// Check that we returned into the original directory after the cleanup
	currWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if currWD != origWD {
		t.Fatalf("Expected to return to the \"%s\" directory after the cleanup. Instead we are in \"%s\"\n", origWD, currWD)
	}
}

func TestTempInitDirRestrictedPermissions(t *testing.T) {

	root, cleanup, err := TempInitDir()
	if err != nil {
		t.Fatal(err)
	}

	d := filepath.Join(root, "d")
	err = os.Mkdir(d, 0700)
	if err != nil {
		t.Fatal(err)
	}

	f := filepath.Join(d, "f")
	err = ioutil.WriteFile(f, []byte{}, 0700)
	if err != nil {
		t.Fatal(err)
	}

	// Set fully restricted permissions on the file and directory
	// so that it is clear the cleanup removes them
	err = os.Chmod(f, 0)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chmod(d, 0)
	if err != nil {
		t.Fatal(err)
	}

	// run the resulting cleaup function for tests below
	cleanup()

	// Check that the testRootDir does not exist after the cleanup
	_, err = os.Stat(root)
	_, ok := err.(*os.PathError)
	if !ok {
		t.Fatalf("Temporary directory \"%s\" remained after cleanup", root)
	}
}

func ExampleTempCloneDir() {

	root, cleanup, err := TempCloneDir("./mock")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestTempCloneDir(t *testing.T) {

	const src string = "./testdata/temp_dir_mocks"

	// Get the values and create the test root dir to be tested
	testRootDir, cleanup, err := TempCloneDir(src)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the returned testRootDir is Stat-able
	testRootInfo, err := os.Stat(testRootDir)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the returned testRootDir is a directory
	if !testRootInfo.IsDir() {
		t.Errorf("Returned temporary path \"%s\" is not a directory", testRootDir)
	}

	diffs, err := TreeDiff(src, testRootDir, ByName, ByDir, BySize, ByPerm, ByTime, ByContent(t))
	if err != nil {
		t.Fatal(err)
	}

	if diffs != nil {
		t.Errorf("Trees at \"%s\" and \"%s\" differ unexpectedly: %v", src, testRootDir, diffs)
	}

	// run the resulting cleaup function for tests below
	cleanup()

	// Check that the testRootDir does not exist after the cleanup
	_, err = os.Stat(testRootDir)
	_, ok := err.(*os.PathError)
	if !ok {
		t.Fatalf("Cloned directory \"%s\" remained after cleanup", testRootDir)
	}
}

func TestTempCloneChdir(t *testing.T) {

	// Capture the old workdir
	origWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	// Get the values and create the test root dir to be tested
	src := "./testdata/temp_dir_mocks"

	// Get the values and clone the test root dir to be tested
	old, cleanup, err := TempCloneChdir(src)
	if err != nil {
		t.Fatal(err)
	}

	if origWD != old {
		t.Fatalf("Got \"%s\" as an old directory instead of the expected \"%s\"\n", old, origWD)
	}

	// Capture the temporary workdir
	tempWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	src = filepath.Join(origWD, src)

	diffs, err := TreeDiff(src, tempWD, ByName, ByDir, BySize, ByPerm, ByTime, ByContent(t))
	if err != nil {
		t.Fatal(err)
	}

	if diffs != nil {
		t.Errorf("Trees at \"%s\" and \"%s\" differ unexpectedly: %v", src, tempWD, diffs)
	}

	// run the resulting cleaup function for tests below
	cleanup()

	// Check that we returned into the original directory after the cleanup
	currWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if currWD != origWD {
		t.Fatalf("Expected to return to the \"%s\" directory after the cleanup. Instead we are in \"%s\"\n", origWD, currWD)
	}
}

func ExampleTempCreateChdir() {

	dirMark := func(fi os.FileInfo) string {
		if fi.IsDir() {
			return "/"
		}
		return ""
	}

	dirs := `
			2001-01-01T01:01:01Z	0750	a/
			2001-01-01T01:01:01Z	0750	b/
			2001-01-01T01:01:01Z	0700	c.txt	"This is a two line\nfile with\ta tab\n"
			2001-01-01T01:01:01Z	0700	d.txt	No need to quote a single line without tabs

			2002-01-01T01:01:01Z	0700	"has\ttab/"
			2002-01-01T01:01:01Z	0700	"has\ttab/e.mb"	"# Markdown...\n\n... also ***possible***\n"

			2002-01-01T01:01:01Z	0700	"\u263asmiles\u263a/"
		`

	nodes, err := TreeParseReader(strings.NewReader(dirs))
	if err != nil {
		log.Fatal(err)
	}

	_, cleanup, err := TempCreateChdir(nodes)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

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
