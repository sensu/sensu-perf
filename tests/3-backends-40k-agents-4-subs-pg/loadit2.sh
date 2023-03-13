#!/usr/bin/env bash

/root/go/bin/loadit -backends ws://backend1:8081,ws://backend2:8081,ws://backend3:8081 \
-subscriptions large_cluster_test_2 \
-count 10000 \
-keepalive-interval 60 \
-keepalive-timeout 360 \
-pprof-port 6061 \
-prom "" \
-base-entity-name "loadit-2-agent"
