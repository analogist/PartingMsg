#!/bin/sh

cd partmsg/
go run create.go &&
cd ..
go build
