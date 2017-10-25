// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fstests // import "didenko.com/go/fstests"

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func ExampleInitTempDir() {

	root, cleanup := InitTempDir()
	if root == "" {
		log.Fatal("Failed to create a temporary directory")
	}
	defer cleanup()

	fmt.Printf("Here goes the code using the temporary directory at %s\n", root)
}

func TestInitTempDir(t *testing.T) {

	// Get the values and create the test root dir to be tested
	testRootDir, cleanup := InitTempDir()
	if testRootDir == "" {
		t.Fatal("Failed to create a temporary directory")
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
