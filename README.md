# Walle Serverless Workflow

## Overview

Walle is an open-source, vendor-neutral and cloud-native serverless workflow management system.

![Workflow](./images/workflow.png)

![Execution](./images/execution.png)

## Quick Start

### Prerequisite

- Kubenetes

- OpenFaaS deployed on your Kubenetes

> [OpenFaaS Deployment guide for Kubernetes](https://docs.openfaas.cOpenfaasom/deployment/kubernetes)

### Deploy Workflow Engine on OpenFaaS

```shell
export OPENFAAS_GATEWAY=XXX
faas-cli deploy -f walle-engine.yml -g $OPENFAAS_GATEWAY
```

### Deploy Other Components on Kubenetes

```shell
chmod +x install.sh
./install.sh
```

### Get UI DashBoard URL

```shell
kubectl get ingress -n walle

NAME              CLASS   HOSTS   ADDRESS        PORTS   AGE
gateway-ingress   nginx   *       192.168.64.2   80      10h
```

Open `ADDRESS` in your browser

## Workflow Example
```yaml
version: 1.0
name: example-workflow
desc: Example Wrokflow
triggers:
- type: http
  name: http-trigger
  async: true
tasks:
# task-1, no dependency
- name: task-1
  type: http
  url: http://example.com/api/task-1
  timeout: 3s
# task-2, no dependency
- name: task-2
  type: http
  url: http://example.com/api/task-2
  timeout: 3s
# task-3, depends on task-1 and task-2, should be run after them
- name: task-3
  type: http
  url: http://example.com/api/task-3
  method: GET
  retry: 3
  timeout: 3s
  depends: [task-1, task-2]
# task-4, only depends on task-1, should be immediately run after task-1
- name: task-4
  type: http
  url: http://example.com/api/task-4
  timeout: 3s
  depends: [task-1]
```

## Architecture

![Architecture Diagram](./images/Architecture.png)