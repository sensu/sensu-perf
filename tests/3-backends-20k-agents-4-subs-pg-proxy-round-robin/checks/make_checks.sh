#!/usr/bin/env bash

for i in {0..1000}; do
    mapfile -t lines < checks.tmpl
    for line in "${lines[@]}"; do
        echo "$line" >> check$i.yml
    done
done
