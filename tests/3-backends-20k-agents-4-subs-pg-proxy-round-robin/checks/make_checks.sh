#!/usr/bin/env bash

for i in {0..1000}; do
    cat << EOF > check$i.yml
type: CheckConfig
api_version: core/v2
metadata:
  name: check`expr $i + 1`
  namespace: default
spec:
  check_hooks: null
  command: '!sensu_test_check!'
  env_vars: null
  handlers: []
  high_flap_threshold: 0
  interval: 10
  low_flap_threshold: 0
  output_metric_format: ""
  output_metric_handlers: null
  proxy_entity_name: ""
  publish: true
  round_robin: true
  proxy_requests:
    entity_attributes:
    - "entity.entity_class == 'proxy'"
    - "entity.labels.proxy_type == 'switch'"
  runtime_assets: null
  stdin: false
  subdue: null
  subscriptions:
  - large_cluster_test_`expr $i % 4 + 1`
  timeout: 0
  ttl: 0
EOF
done
