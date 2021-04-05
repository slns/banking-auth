#!/bin/bash
SERVER_ADDRESS=localhost \
SERVER_PORT=8181 \
DB_USER=root \
DB_PASSWD=root \
DB_ADDR=localhost \
DB_PORT=8889 \
DB_NAME=banking \
go run main.go
