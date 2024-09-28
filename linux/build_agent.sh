#!/bin/bash

cd agent

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o agent server.go

# Move the built binary to the assets directory
mv agent ../assets/

echo "---- Agent: build completed and moved to ./assets/"
