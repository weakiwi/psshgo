#!/bin/bash -x

#go get -u github.com/weakiwi/gosshtool
#go get -u github.com/urfave/cli
go build
mv psshgo release/
