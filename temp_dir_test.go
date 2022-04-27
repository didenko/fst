// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func ExampleTempInitDir() {
	lg := log.New(os.Stderr, "ExampleTempInitDir", log.LUTC|log.Ldate|log.Ltime)
	root, cleanup := TempInitDir(lg)
	defer cleanup()
	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestTempInitDir(t *testing.T) {

	// Get the values and create the test root dir to be tested
	testRootDir, cleanup := TempInitDir(t)

	// Check that the returned testRootDir is Stat-able
	testRootInfo, err := os.Stat(testRootDir)
	if err != nil {
		t.Fatalf("While learning about the temporary directory %q: %s", testRootDir, err)
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
	old, cleanup := TempInitChdir(t)

	if origWD != old {
		t.Fatalf("Got \"%s\" as an old directory instead of the expected \"%s\"\n", old, origWD)
	}

	// Capture the temporary workdir
	tempWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Workind directory is indetermined: %s", err)
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
		t.Fatalf("Learning about the original directory after the tests: %s", err)
	}

	if currWD != origWD {
		t.Fatalf("Expected to return to the \"%s\" directory after the cleanup. Instead we are in \"%s\"\n", origWD, currWD)
	}
}

func TestTempInitDirRestrictedPermissions(t *testing.T) {

	root, cleanup := TempInitDir(t)

	d := filepath.Join(root, "d")
	err := os.Mkdir(d, 0700)
	if err != nil {
		t.Fatalf("Creating a test directory %q: %s", d, err)
	}

	f := filepath.Join(d, "f")
	err = ioutil.WriteFile(f, []byte{}, 0700)
	if err != nil {
		t.Fatalf("Creating a test directory %q: %s", f, err)
	}

	// Set fully restricted permissions on the file and directory
	// so that it is clear that the cleanup removes them
	err = os.Chmod(f, 0)
	if err != nil {
		t.Fatalf("Permissions on a test directory %q: %s", f, err)
	}

	err = os.Chmod(d, 0)
	if err != nil {
		t.Fatalf("Permissions on a test directory %q: %s", d, err)
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
	lg := log.New(os.Stderr, "ExampleTempCloneDir", log.LUTC|log.Ldate|log.Ltime)
	root, cleanup := TempCloneDir(lg, "./mock")
	defer cleanup()
	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestTempCloneDir(t *testing.T) {

	const src string = "./testdata/temp_dir_mocks"

	// Get the values and create the test root dir to be tested
	testRootDir, cleanup := TempCloneDir(t, src)

	// Check that the returned testRootDir is Stat-able
	testRootInfo, err := os.Stat(testRootDir)
	if err != nil {
		t.Fatal(err)
	}

	// Check that the returned testRootDir is a directory
	if !testRootInfo.IsDir() {
		t.Errorf("Returned temporary path \"%s\" is not a directory", testRootDir)
	}

	diffs := TreeDiff(t, src, testRootDir, ByName, ByDir, BySize, ByPerm, ByTime, ByContent(t))

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
		t.Fatalf("Original working directory is indetermined: %s", err)
	}

	// Get the values and create the test root dir to be tested
	src := "./testdata/temp_dir_mocks"

	// Get the values and clone the test root dir to be tested
	old, cleanup := TempCloneChdir(t, src)

	if origWD != old {
		t.Fatalf("Got \"%s\" as an old directory instead of the expected \"%s\"\n", old, origWD)
	}

	// Capture the temporary workdir
	tempWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Temporary working directory is indetermined: %s", err)
	}

	src = filepath.Join(origWD, src)

	diffs := TreeDiff(t, src, tempWD, ByName, ByDir, BySize, ByPerm, ByTime, ByContent(t))

	if diffs != nil {
		t.Errorf("Trees at \"%s\" and \"%s\" differ unexpectedly: %v", src, tempWD, diffs)
	}

	// run the resulting cleaup function for tests below
	cleanup()

	// Check that we returned into the original directory after the cleanup
	currWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("Original working directory is indetermined after cleanup: %s", err)
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

	lg := log.New(os.Stderr, "ExampleTempCreateChdir", log.LUTC|log.Ldate|log.Ltime)

	nodes := []*Node{
		{0750, Rfc3339(lg, "2001-01-01T01:01:01Z"), "a/", ""},
		{0750, Rfc3339(lg, "2001-01-01T01:01:01Z"), "b/", ""},
		{0700, Rfc3339(lg, "2001-01-01T01:01:01Z"), "c.txt", "This is a two line\nfile with\ta tab\n"},
		{0700, Rfc3339(lg, "2001-01-01T01:01:01Z"), "d.txt", "A single line without tabs"},
		{0700, Rfc3339(lg, "2002-01-01T01:01:01Z"), "has\ttab/", ""},
		{0700, Rfc3339(lg, "2002-01-01T01:01:01Z"), "has\ttab/e.mb", "# Markdown...\n\n... also ***possible***\n"},
		{0700, Rfc3339(lg, "2002-01-01T01:01:01Z"), "\u263asmiles\u263a/", ""},
	}

	_, cleanup := TempCreateChdir(lg, nodes)
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
