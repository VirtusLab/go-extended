# go-extended

[![Version](https://img.shields.io/badge/version-v0.0.1-brightgreen.svg)](https://github.com/VirtusLab/go-extended/releases/tag/v0.0.1)
[![Travis CI](https://img.shields.io/travis/VirtusLab/go-extended.svg)](https://travis-ci.org/VirtusLab/go-extended)
[![Go Report Card](https://goreportcard.com/badge/github.com/VirtusLab/go-extended "Go Report Card")](https://goreportcard.com/report/github.com/VirtusLab/go-extended)
[![GoDoc](https://godoc.org/github.com/VirtusLab/go-extended?status.svg "GoDoc Documentation")](https://godoc.org/github.com/VirtusLab/go-extended/)

Things missing or not belonging in the standard go library

**No external dependencies**, with two exceptions:
- go standard library
- test libraries

* [Installation](README.md#installation)
  * [Via Go](README.md#via-go)
* [Usage](README.md#usage)
* [Contribution](README.md#contribution)
* [Development](README.md#development)
* [The Name](README.md#the-name)


## Installation
#### Via Go

```console
$ go get github.com/VirtusLab/go-extended
```

## Usage

See [GoDoc Documentation](https://godoc.org/github.com/VirtusLab/go-extended/)
and [tests](https://github.com/VirtusLab/render/blob/master/renderer/render_test.go) 
for usage examples.

## Contribution

Feel free to file [issues](https://github.com/VirtusLab/go-extended/issues) 
or [pull requests](https://github.com/VirtusLab/go-extended/pulls).

## Development

    export GOPATH=$HOME/go
    export PATH=$PATH:$GOPATH/bin
    
    mkdir -p $GOPATH/src/github.com/VirtusLab
    cd $GOPATH/src/github.com/VirtusLab/go-extended
    git clone git@github.com:VirtusLab/go-extended.git
    
    go get -u github.com/golang/dep/cmd/dep
    make all

## The name

We believe in obvious names. It extends go. It's `go-extended`.
