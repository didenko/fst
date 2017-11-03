# FileSystem TESTS
[![GoDoc](https://godoc.org/go.didenko.com/fstest?status.svg)](https://godoc.org/go.didenko.com/fstest)
[![Build Status](https://travis-ci.org/didenko/fstest.svg?branch=master)](https://travis-ci.org/didenko/fstest)

## Purpose

It is at times desireable to test  a program behavior which creates or modifies files and directories. Such tests may be reasonable involved especially if they involve checking of attributes like ownership, permissions, or timestamps. It becomes an extra burden as such filesystem manipulation has to be tested itself - so one ends up with tests of tests.

The `fstest` library is a tested set of functions aiming to help with it. It makes simple to clone and compare filesystem trees for testing puposes.

## Limitations

Functions in `fstest` expect a reasonably shallow and small directory structures to deal with, as that is what usually happens in testing. In particular during build-up, tear-down, and comparisons it creates collections of filesystem object names in memory. It is not nesessarily efficient in that sense, but allows for more graceful permissions handling.

If you are concerned, that you will hold a few copies of full filenames' lists during the execution, then this library may be a poor match to your needs.
