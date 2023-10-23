#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: $0 <log-filename>"
    exit 1
fi

for ((i=0; i<4; i++)); do
    toAdd=$(cat "$1" | grep "$i]: | janitor: pindexes to add:" | awk '{sum += $9} END {print sum}')
    toRemove=$(cat "$1" | grep "$i]: | janitor: pindexes to remove:" | awk '{sum += $9} END {print sum}')
    result=$(echo "add:$toAdd - remove:$toRemove")
    re='^[0-9]+$'
    if ! [[ $toAdd =~ $re ]] || ! [[ $toRemove =~ $re ]] ;
    then
        echo "ERR: Node[$i]: toAdd: $toAdd, toRemove: $toRemove"
    else
        echo "Node[$i]: $result = $(($toAdd - $toRemove))"
    fi
done