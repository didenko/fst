// Copyright 2017-2019 Vlad Didenko. All rights reserved.
// See the included LICENSE.md file for licensing information

// Package fst is a collection of functions to help
// testing filesyste objects modifications. It focuses
// on creating and cleaning up a baseline filesyetem state.
//
// The following are addressed use cases:
//
// 1. Create a directory hierarchy via an API
//
// 2. Create a directory hierarchy via a copy of a template
//
// 3. Write a provided test mock data to files
//
// 4. Contain all test activity in a temporatry directory
//
// 5. Compare two directories recursively
package fst // import "go.didenko.com/fst"
