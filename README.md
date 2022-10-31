# Walle Serverless Workflow

## Overview

Walle is an open-source, vendor-neutral and cloud-native serverless workflow management system.

![New Workflow](./images/New%20Workflow.png)

![Workflow](./images/workflow.png)

![Execution](./images/execution.png)

## Architecture

### Key Components

![Architecture Diagram](./images/Architecture.png)

### Workflow Engine

![Workflow Engine Diagram](./images/Workflow%20Engine.png)

### Workflow Executor
![Workflow Executor Diagram](./images/Workflow%20Executor.png)

## Parallelism

### Sync Diagram

![Sync Diagram](./images/Sync.png)

### Batch Async

![Batch Async Diagram](./images/Batch%20Async.png)

### 🌟 Full Async
![Full Async Diagram](./images/Fully%20Async.png)

## Workflow Pattern

### Sequential

![Sequential Diagram](./images/Seq.png)

### Parallel

![Parallel Diagram](./images/Para.png)

### Synchronize

![Synchronize Diagram](./images/Sync%20Pattern.png)

### Merge

![Merge Diagram](./images/Merge.png)

## Quick Start

### Prerequisite

- Kubernetes

- OpenFaaS deployed on your Kubernetes

> [OpenFaaS Deployment guide for Kubernetes](https://docs.openfaas.cOpenfaasom/deployment/kubernetes)

### Deploy Workflow Engine on OpenFaaS

```shell
export OPENFAAS_GATEWAY=XXX
faas-cli deploy -f walle-engine.yml -g $OPENFAAS_GATEWAY
```

### Deploy Other Components on Kubernetes

```shell
chmod +x install.sh
./install.sh
```

### Get UI Dashboard URL

```shell
kubectl get ingress -n walle

NAME              CLASS   HOSTS   ADDRESS        PORTS   AGE
gateway-ingress   nginx   *       192.168.64.2   80      10h
```

Open `ADDRESS` in your browser

## Workflow Example
```yaml
name: example-workflow
triggers:
- type: http
  name: http-trigger
  async: true
tasks:
- name: task-1 # task name
  type: http # task type
  url: http://example.com/api/task-1 # the cloud function URL
  method: GET # HTTP method
  headers: # HTTP headers
  - name: Content-Type
    value: application/json
  body: {"message": "Hello Walle!"} # HTTP body
  retry: 3 # max retry times if the task failed
  timeout: 10s # max request timeout
  depends: [] # dependency list
- name: task-2 # Task 2, no dependency
  url: http://example.com/api/task-2
- name: task-3
  url: http://example.com/api/task-3
  depends: [task-1, task-2] # depends on both Task 1 and Task 2, should be run after them
- name: task-4
  url: http://example.com/api/task-4
  depends: [task-1] # Only depends on Task 1 and should be immediately run after Task 1
```

## Support Workflow Patterns

```yaml
name: sequence-pattern
tasks:
- name: task-1
  url: http://example.com/api/task-1
- name: task-2
  url: http://example.com/api/task-2
  depends: [task-1]
- name: task-3
  url: http://example.com/api/task-3
  depends: [task-2]
```

```yaml
name: parallel-split-pattern
tasks:
- name: task-1
  url: http://example.com/api/task-1
- name: task-2
  url: http://example.com/api/task-2
  depends: [task-1]
- name: task-3
  url: http://example.com/api/task-3
  depends: [task-1]
```

```yaml
name: synchronization-pattern
tasks:
- name: task-1
  url: http://example.com/api/task-1
- name: task-2
  url: http://example.com/api/task-2
- name: task-3
  url: http://example.com/api/task-3
  depends: [task-1, task-2]
```