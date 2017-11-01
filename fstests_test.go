// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstests // import "go.didenko.com/fstests"

import (
	"fmt"
	"log"
	"os"
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

func ExampleCloneTempDir() {

	root, cleanup, err := CloneTempDir("./mock")
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestCloneTempDir(t *testing.T) {

	const src string = "./mock"

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

	if diffs := TreeDiffs(src, testRootDir); diffs != nil {
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
