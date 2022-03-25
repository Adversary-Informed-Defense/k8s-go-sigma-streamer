# Sigma Engine

Sigma engine for pre-SIEM signature-based detection of logs.

## Features

- Analyze streaming logs using Sigma rules
- Kafka ingest and and analysis output
- Create and test Sigma rules using a variety of built-in sample logs
- Load-based autoscaling using Kubernetes Horizontal Pod Autoscalers (HPA)
- Written in GoLang

## Getting Started

### Install Dependencies

- Docker <https://docs.docker.com/get-docker/>
- Docker Compose <https://docs.docker.com/compose/install/>
- Kubernetes + kubectl (Kubernetes through Docker Desktop on MacOS is preferred)
- Helm  <https://helm.sh/docs/intro/install/>
- Helmfile <https://github.com/roboll/helmfile/releases>
- Kubernetes Metrics API

  ```bash
  kubectl apply -f ./deploy/resources/metrics_server.yaml
  ```

### Configure Rule Sources, Build, Run, and Test

To start out, run the following in order:

Edit the rules `helpers/rule_sources.json` (must have access to repos specified)

```bash
make build   # Build images from sigma engine and integration test
make deploy  # Deploy k8s cluster
make test    # Run integration test to run sample data through engine
./helpers/consume.sh # Consume logs from output topic "alerts"
```

### Disclaimer
This application uses the go-sigma-rule-engine library (https://github.com/markuskont/go-sigma-rule-engine) available under the Apache 2.0 license, which can be obtained from http://www.apache.org/licenses/LICENSE-2.0.
