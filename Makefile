SHELL:=/bin/bash

BIN_NAME=go-email-mock
PID_FILE=run.pid

build:
	go build -o ${BIN_NAME} ./...

start: build
	./${BIN_NAME} & echo $$! > ${PID_FILE}

stop:
	kill $$(cat ${PID_FILE})
	rm ${PID_FILE}

test:
	./run-tests.sh
