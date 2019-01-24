// Copyright 2017 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

package fst // import "go.didenko.com/fst"

import (
	"testing"
)

func TestByContent(t *testing.T) {

	files := []*Node{
		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "1_same_a", "a 1 b 2 c 3 d 4\n"},
		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "1_same_b", "a 1 b 2 c 3 d 4\n"},

		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "2_diff_a", "a 1 b 2 c 3 d 4\n"},
		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "2_diff_b", "a 1 b 2 b 2\n"},

		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "3_empty_a", ""},
		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "3_empty_b", ""},

		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "4_length_a", "a 1 b 2 c 3 d 4\n"},
		&Node{0700, Rfc3339(t, "2017-11-08T23:11:17Z"), "4_length_b", "a 1 b 2 c"},
	}

	_, cleanup, err := TempCreateChdir(files)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()

	less := ByContent(t)

	fips, err := MakeFipSlice(
		"1_same_a", "1_same_b",
		"2_diff_a", "2_diff_b",
		"3_empty_a", "3_empty_b",
		"4_length_a", "4_length_b",
	)
	if err != nil {
		t.Fatal(err)
	}

	if less(fips[0], fips[1]) {
		t.Errorf("Files %v and %v with same content ranked as ordered\n", fips[0].Name(), fips[1].Name())
	}

	if less(fips[2], fips[3]) {
		t.Errorf("File %v is incorrectly ranked less than %v\n", fips[2].Name(), fips[3].Name())
	}

	if !less(fips[3], fips[2]) {
		t.Errorf("File %v is not ranked less than %v as expected\n", fips[3].Name(), fips[2].Name())
	}

	if less(fips[4], fips[5]) {
		t.Errorf("Empty files %v and %v incorrectly ranked as ordered\n", fips[4].Name(), fips[5].Name())
	}

	if less(fips[6], fips[7]) {
		t.Errorf("File %v is incorrectly ranked less than %v\n", fips[6].Name(), fips[7].Name())
	}

	if !less(fips[7], fips[6]) {
		t.Errorf("File %v is not ranked less than %v as expected\n", fips[7].Name(), fips[6].Name())
	}
}
