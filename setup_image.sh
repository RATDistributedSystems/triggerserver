#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -a --installsuffix cgo --ldflags="-s" -o triggerserver

# Build the image
docker build -t rattrigger .

# Remove remnants
rm -f triggerserver