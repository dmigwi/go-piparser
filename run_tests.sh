#!/bin/bash

# Check for internet connectivity to decide if all tests will be run.
echo -e "GET http://github.com HTTP/1.0\n\n" | nc github.com 80 > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "System detected internet connectivity"
    go test -tags fulltest ./... -v 
else
    echo "Ping to github.com failed"
    go test ./... -v 
fi