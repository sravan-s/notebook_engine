#!/bin/bash

cd ../agent
go build ./server.go
mv server ../linux/agent/agent