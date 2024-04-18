#!/bin/bash

# Create a test binary which will be used to run each test individually
go test -c -o tests.test

# Run each test individually, printing "." for successful tests, or the test name
for test in $(go test -list . | grep -E "^(Test|Example)"); do
    ./tests.test -test.run "^$test$" &>/dev/null
    if [ $? -eq 0 ]; then
        echo -n "."
    else
        echo -e "\n$test failed"
    fi
done

echo ""
