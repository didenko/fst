// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstest // import "go.didenko.com/fstest"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func ExampleInitTempDir() {

	root, cleanup, err := InitTempDir()
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestInitTempDir(t *testing.T) {

	// Get the values and create the test root dir to be tested
	testRootDir, cleanup, err := InitTempDir()
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
	if !ok {
		t.Fatalf("Temporary directory \"%s\" remained after cleanup", testRootDir)
	}
}

func TestInitTempDirRestrictedPermissions(t *testing.T) {

	root, cleanup, err := InitTempDir()
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

func ExampleCloneTempDir() {

	root, cleanup, err := CloneTempDir("./mock")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestCloneTempDir(t *testing.T) {

	const src string = "./temp_dir_mocks"

	// Get the values and create the test root dir to be tested
	testRootDir, cleanup, err := CloneTempDir(src)
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

	if diffs := TreeDiff(src, testRootDir, ByName, ByDir, BySize, ByPerm, ByTime); diffs != nil {
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
