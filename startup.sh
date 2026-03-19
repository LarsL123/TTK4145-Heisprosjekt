#!/bin/bash
# startup.sh

go build -o elevator src/main.go


while true;do
    ./elevator "$1"
    echo "Died, restarting in 200ms..."
    sleep 0.2
done