#!/bin/bash

./main &
sleep 1
echo
echo -n "alex" | nc -i 1 localhost 8001
kill $!
