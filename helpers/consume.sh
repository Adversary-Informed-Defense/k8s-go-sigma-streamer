#!/bin/bash
kubectl exec -it kafka-0 -- /bin/bash \
  -c 'unset JMX_PORT; /opt/bitnami/kafka/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic alerts --group=unique'
