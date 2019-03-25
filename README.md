# go-extended

[![Version](https://img.shields.io/badge/version-v0.0.4-brightgreen.svg)](https://github.com/VirtusLab/go-extended/releases/tag/v0.0.4)
[![Travis CI](https://img.shields.io/travis/VirtusLab/go-extended.svg)](https://travis-ci.org/VirtusLab/go-extended)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/go-extended "Go Report Card")](https://goreportcard.com/report/github.com/VirtusLab/go-extended)
[![GoDoc](https://godoc.org/github.com/VirtusLab/go-extended?status.svg "GoDoc Documentation")](https://godoc.org/github.com/VirtusLab/go-extended/)

* [Goals](README.md#goals)
* [Installation](README.md#installation)
  * [Via Go](README.md#via-go)
* [Usage](README.md#usage)
* [Contribution](README.md#contribution)
* [Development](README.md#development)
* [The Name](README.md#the-name)

## Goals

Things missing or not belonging in the standard go library

**No external dependencies**, with two exceptions:
- go standard library
- tests

## Installation
#### Via Go

```console
$ go get github.com/VirtusLab/go-extended
```

## Usage

See [GoDoc Documentation](https://godoc.org/github.com/VirtusLab/go-extended/)
and the tests, [e.g.](https://github.com/VirtusLab/go-extended/blob/master/pkg/renderer/render_test.go) 
for usage examples.

### Notable features

- simple [`renderer`](https://godoc.org/github.com/VirtusLab/go-extended/pkg/renderer) that extends [`text/template`](https://golang.org/pkg/text/template/)
- easy to use [`matcher`](https://godoc.org/github.com/VirtusLab/go-extended/pkg/matcher) that extends [`regexp`](https://golang.org/pkg/regexp/)

## Contribution

Feel free to file [issues](https://github.com/VirtusLab/go-extended/issues) 
or [pull requests](https://github.com/VirtusLab/go-extended/pulls).

## Development

    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    mkdir -p $GOPATH/src/github.com/VirtusLab
    cd $GOPATH/src/github.com/VirtusLab
    git clone git@github.com:VirtusLab/go-extended.git
    cd go-extended
    
    go get -u github.com/golang/dep/cmd/dep
    make all

## The name

We believe in obvious names. It extends go. It's `go-extended`.
