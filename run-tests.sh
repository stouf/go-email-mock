#!/usr/bin/env bash

set +e
set +u

make start
sleep 1
go test ./...
make stop
