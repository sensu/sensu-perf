#!/usr/bin/env bash

/root/go/bin/loadit -backends ws://backend1:8081,ws://backend2:8081,ws://backend3:8081 \
-subscriptions large_cluster_test_1 \
-count 2000 \
-keepalive-interval 60 \
-keepalive-timeout 360 \
-pprof-port 6060 \
-prom "" \
-base-entity-name "loadit-1-agent"
