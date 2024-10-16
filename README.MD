# StackState Kubernetes Extension Integration

[StackState](https://stackstate.com) is full-stack observability platform with tons of extension and integration possibilities.

Out of the box StackState supports majority of Kubernetes resources, but not all.  This integration attempts to fill the
gaps with the default StackState implementation.  Overtime more functionality will be added as needed.

Additional resources supported,

- storage.k8s.io/v1/StorageClass

## Helm Deployment

Setup a basic values.yaml containing connectivity information about your StackState instance. 

```yaml
serverUrl:  https://xxxx.    # Suse Observability Url
apiKey: xxx                  # Suse Observability Api Key
clusterName: lab             # Cluster name as defined for the Kubernetes StackPack instance in Suse Observability
```

Then run the following Helm commands,

```bash
helm repo add sts-ext https://ravan.github.io/helm-charts
helm upgrade --install --namespace stackstate-extensions --create-namespace -f values.yaml sts sts-ext/stackstate-k8s-ext                            

```

## Configuration

### Environment variables

Environment variables can be used to set up the scanner.

| Variable                     | Required | Default        | Description                                      |
|------------------------------|----------|----------------|--------------------------------------------------|
| STACKSTATE_API_URL           | yes      |                | The url to the StackState server                 |
| STACKSTATE_API_KEY           | yes      |                | The api key used by the StackState server        |
| KUBERNETES_KUBECONFIG        | no       | ~/.kube/config | KubeConfig file to use to connect to K8s         |
| KUBERNETES_IN_CLUSTER        | no       | false          | Use in cluster K8s connection                    |
| KUBERNETES_CLUSTER           | yes      |                | Kubernetes Cluster name as defined in StackState |
| INSTANCE_TYPE                | yes      | k8s-ext        | The Custom Syn StackPack instance type.          |
| INSTANCE_URL                 | yes      |                | The Custom Syn StackPack instance url            |


### Config File

Configuration could also be in a yaml file.
The location of the config file can be set in the `CONFIG_FILE`environment variable

```yaml
stackstate:
  api_url: "https://xxx.stackstate.io"
  api_key: "5385xxxx"
  
kubernetes:
  cluster: retailstore
  in_cluster: false
  kubeconfig: ./kubeconfig

instance:
  type: crd
  url: retailstore

```


## Development

### Prerequisites

- [Taskfile](https://taskfile.dev/installation/)


### Setup environment variables

Set up a `.env` file for configuring the env.

```
USER="<my dockerhub user name>"
PWD="<my dockerhub password>"
STS_URL=https://myinstance.stackstate.io
STS_API_KEY=xxxx
```
### Build

```shell
task docker-build
```

### Push

```shell
task docker-push
```

