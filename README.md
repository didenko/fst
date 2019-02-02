# **fst** - File System Testing aids

[![GoDoc](https://godoc.org/go.didenko.com/fst?status.svg)](https://godoc.org/go.didenko.com/fst)
[![Build Status](https://travis-ci.org/didenko/fst.svg?branch=master)](https://travis-ci.org/didenko/fst)
[![Go Report Card](https://goreportcard.com/badge/go.didenko.com/fst)](https://goreportcard.com/report/go.didenko.com/fst)

The suggested package name pronounciation is _"fist"_.

> This is the version 1 branch of the `fst` package. The upcoming version 2 simplifies the API significantly. It is, however, backward incompatible and will require some client code modifications. The version 2 is in beta until about the end of Martch 2019 - see the v2 branch. Please, [file an issue](https://github.com/didenko/fst/issues) with a feedback or a bug.

## Purpose

At times it is desireable to test a program behavior which creates or modifies files and directories. Such tests may be quite involved especially if checking permissions or timestamps. A proper cleanup is also considered a nuisance. The whole effort becomes extra burdesome as such filesystem manipulation has to be tested itself - so one ends up with tests of tests.

The `fst` library is a tested set of functions aiming to alleviate the burden. It makes creating and comparing filesystem trees of regular files and directories for testing puposes simple.

## Highlights

The three most used functtions in the `fst` library are [_TempCloneChdir_](#TempCloneChdir), [_TempCreateChdir_](#TempCreateChdir), and [_TreeDiff_](#TreeDiff). For details on these and other functions, see the examples and documentation at the https://godoc.org/go.didenko.com/fst page.

### <span id="TempCloneChdir" />[TempCloneChdir](https://godoc.org/go.didenko.com/fst#TempCloneDir)

It is used to clone an existing directory with all it's content, permissions, and timestamps. Consider this example:

```go
old, cleanup, err := TempCloneChdir("mock")
if err != nil {
  t.Fatal(err)
}
defer cleanup()
```

If an error was returned, then no temporary directory was left behind (after a reasonable cleanup effort). Otherwise, the _mock_ directory's content will be cloned into a new temporary directory and the calling process will change into it. The _old_ variable in the example will contain the original directory where the running process was at the time of the _TempCloneChdir_ call.

The ***cleanup*** function has the code to change back to the original directory and then delete the temporary directory.

As the _TempCloneChdir_ relies on the `TreeCopy` function, it will attempt to recreate both permissions and timestamps from the source directory. Keep in mind, that popular version control systems like _Git_ and _Mercurial_ do not preserve original files' timestamps. If your tests rely on timestamped files or directories then _TreeCreate_ or it's derivative _TempCreateChdir_ functions are your friends.

### <span id="TempCreateChdir" />[TempCreateChdir](https://godoc.org/go.didenko.com/fst#TempCreateChdir)

The _TempCreateChdir_ function provides an API-like way to create and populate a temporary directory tree for testing. It takes an _io.Reader_, from which is expects to receive lines with ***tab-separated*** fields describing the directories and files to be populated. Here is an example:

```go
tree := `
2017-11-12T13:14:15Z	0750	settings/
2017-11-12T13:14:15Z	0640	settings/theme1.toml	key = val1
2017-11-12T13:14:15Z	0640	settings/theme2.toml	key = val2
`
treeR := strings.NewReader(tree)
old, cleanup, err = TempCreateChdir(treeR)
if err != nil {
  t.Fatal(err)
}
defer cleanup()
```

If there is no error, _TempCreateChdir_ will create the `settings` directory  with files `theme1.toml` and `theme2.toml`, with the specified key/value pairs as content in a new directory.

The _TempCreateChdir_ function removes the temporary directory it creates as a part of the cleanup logic. It also does a best-effort attempt to remove the directory is an error occurred during it's operation.

### <span id="TreeDiff" />[_TreeDiff_](https://godoc.org/go.didenko.com/fst#TreeDiff)

The _TreeDiff_ function produces a human-readable output of differences between two directory trees for diagnostic purposes. The resulting slice of strings is empty if no differences are found.

Criteria for comparing filesystem objects varies based on a task, so _TreeDiff_ takes a list of comparator functions. The most common ones are provided with the `fst` package. Users are free to provide their own additional comparators which satisfy the [_FileRank_](https://godoc.org/go.didenko.com/fst#FileRank) signature.

A quick example of a common _TreeDiff_ use:

```go
diffs, err := TreeDiff(
  "dir1", "dir2",
  ByName, ByDir, BySize, ByContent(t))

if err != nil {
  t.Fatal(err)
}

if diffs != nil {
  t.Logf("Differences between dir1 and dir2:\n%v\n", diffs)
}
```

Note, that while the _BySize_ comparator is redundant in presense of the _ByContent_ comparator, in most cases the cheaper size comparison will avoid a more expensive content comparison. The comparator order is significant, because once an earlier comparator returns `true`, the later comparators are not run.

It is easy to provide overly restrictive permissions using the tree cloning and tree creation functions. When unable to access needed information, _TreeDiff_ will return a related error. While specifics may vary it is often safest to set user read and execute permissions for directories and user read permission for files.

## Limitations

Functions in `fst` expect a reasonably shallow and small directory structures to deal with, as that is what usually happens in testing. During build-up, tear-down, and comparisons it creates collections of filesystem object names in memory. It is not necessarily efficient but allows for more graceful permissions handling.

If you are concerned that you will hold a few copies of full filenames' lists during the execution, then this library may be a poor match to your needs.
