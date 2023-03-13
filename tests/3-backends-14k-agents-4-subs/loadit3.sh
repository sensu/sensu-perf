#!/usr/bin/env bash

/root/go/bin/loadit -backends ws://backend1:8081,ws://backend2:8081,ws://backend3:8081 \
-subscriptions large_cluster_test_3 \
-count 3500 \
-keepalive-interval 60 \
-keepalive-timeout 360 \
-pprof-port 6062 \
-prom "" \
-base-entity-name "loadit-3-agent"
