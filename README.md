# node-zone-controller

A Kubernetes controller that automatically applies zone labels to nodes based on custom label.

## Overview

This controller watches for nodes with the `altinity.cloud/auto-zone` label and automatically applies the standard Kubernetes zone labels to those nodes. This is useful for environments where nodes need custom zone assignment that differs from the cloud provider's default zone detection.

## How it Works

The controller continuously monitors all nodes in the cluster. When it detects a node with the configured label key (default: `altinity.cloud/auto-zone`), it automatically applies both current and legacy Kubernetes zone labels with the specified value.

## Configuration

### Label Key
The label key to watch can be configured via:
- **CLI flag**: `--label-key="your.label/key"`
- **Environment variable**: `LABEL_KEY=your.label/key`
- **Default**: `altinity.cloud/auto-zone`

Note: CLI flag takes precedence over environment variable.

## Supported Labels

### Input Label
- Configurable label key (default: `altinity.cloud/auto-zone`): The desired zone value for the node

### Applied Labels
When a node has the configured input label, the controller will automatically apply:
- `topology.kubernetes.io/zone`: The current standard Kubernetes zone label
- `failure-domain.beta.kubernetes.io/zone`: The legacy zone label for backward compatibility

## Examples

### Using default label key
To assign a node to zone `eu-central-1a`:

```bash
kubectl label node <node-name> altinity.cloud/auto-zone=eu-central-1a
```

### Using custom label key
To use a custom label key via environment variable:

```bash
export LABEL_KEY=mycompany.io/zone
# Then start the controller
```

Or via CLI flag:

```bash
./controller --label-key="mycompany.io/zone"
```

Then label your nodes:

```bash
kubectl label node <node-name> mycompany.io/zone=us-west-2a
```

The controller will automatically apply:
- `topology.kubernetes.io/zone=us-west-2a`
- `failure-domain.beta.kubernetes.io/zone=us-west-2a`

## Deployment

Deploy the controller using the provided manifests (replace latest tag for production and use IfNotPresent pull policy):

```bash
kubectl apply -f deploy/rbac.yaml
VERSION=latest envsubst < deploy/deployment.yaml | kubectl apply -f -
```

This will create:
- A deployment running the controller
- RBAC permissions for the controller to read and patch nodes

