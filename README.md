# **fst** - File System Testing aids

[![GoDoc](https://godoc.org/go.didenko.com/fst?status.svg)](https://godoc.org/go.didenko.com/fst)
[![Build Status](https://travis-ci.org/didenko/fst.svg?branch=master)](https://travis-ci.org/didenko/fst)
[![Go Report Card](https://goreportcard.com/badge/go.didenko.com/fst)](https://goreportcard.com/report/go.didenko.com/fst)

The package name should be pronounced as _"fist"_.

## Purpose

At times it is desireable to test a program behavior which creates or modifies files and directories. Such tests may be quite involved especially if checking permissions or timestamps. It becomes an extra burden as such filesystem manipulation has to be tested itself - so one ends up with tests of tests.

The `fst` library is a tested set of functions aiming to help with it. It makes cloning and comparing filesystem trees for testing puposes simple.

## Limitations

Functions in `fst` expect a reasonably shallow and small directory structures to deal with, as that is what usually happens in testing. During build-up, tear-down, and comparisons it creates collections of filesystem object names in memory. It is not nesessarily efficient but allows for more graceful permissions handling.

If you are concerned that you will hold a few copies of full filenames' lists during the execution, then this library may be a poor match to your needs.
