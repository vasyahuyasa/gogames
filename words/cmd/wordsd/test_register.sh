#!/bin/bash

while read line
do
    echo -n $line | nc -uc localhost 48879
done <<EOF
{"command": "register", "name": "test1"}
{"command": "register", "name": "test2"}
{"command": "register", "name": "test3"}
EOF
