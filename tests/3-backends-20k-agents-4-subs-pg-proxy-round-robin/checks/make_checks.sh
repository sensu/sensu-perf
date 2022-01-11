#!/usr/bin/env bash

for i in {0..1000}; do
    while read line
    do
        eval echo "$line" >> check$i.yml
    done < "./checks.tmpl"
done
