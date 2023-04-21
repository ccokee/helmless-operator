# Helmless Operator

The Helmless Operator is a Kubernetes operator designed to deploy any Helm chart using values from a GitHub Gist. It allows users to create a custom resource that specifies the chart repository, chart name, chart version, namespace, and the public Gist containing the values file.

## Prerequisites

- Kubernetes cluster
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) installed and configured to interact with the cluster
- [Helm](https://helm.sh/docs/intro/install/) v3.x installed

## Installation

1. Clone the Helmless Operator repository:

```sh
git clone https://github.com/ccokee/helmless-operator.git
cd helmless-operator
```

2. Deploy the Helmless Operator CRD:

```sh
kubectl apply -f config/crd/bases/cache.redrvm.cloud_helmlesses.yaml
```

3. Install the operator's deployment and related resources:

```sh
kubectl apply -f helmless-operator.yaml
```

## Usage

1. Create a `helmless-deployer.yaml` file containing the HelmLess custom resource:

```yaml
apiVersion: cache.redrvm.cloud/v1alpha1
kind: HelmLess
metadata:
  labels:
    app.kubernetes.io/name: helmless
    app.kubernetes.io/instance: helmless-0
    app.kubernetes.io/part-of: helmless
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: helmless-operator
  name: helmless-deployer
spec:
  chartRepo: <your-chart-repo>
  chartName: <your-chart-name>
  chartVersion: "<your-chart-version>"
  namespace: default
  publicGist: https://gist.github.com/username/gist_id
```

Replace `<your-chart-repo>`, `<your-chart-name>`, `<your-chart-version>` with the corresponding details of the Helm chart you want to deploy.

2. Apply the HelmLess custom resource:

```sh
kubectl apply -f helmless-deployer.yaml
```

The Helmless Operator will now deploy the specified Helm chart using the values provided in the public Gist.

## Example

To deploy any Helm chart, follow these steps:

1. Create a public Gist containing your desired values for the chart. The Gist URL will look like `https://gist.github.com/username/gist_id`.

2. Update the `publicGist` field in the `helmless-deployer.yaml` file with your Gist URL. Also, update the `chartRepo`, `chartName`, and `chartVersion` fields with the details of the Helm chart you want to deploy.

3. Apply the HelmLess custom resource:

```sh
kubectl apply -f helmless-deployer.yaml
```

The Helmless Operator will deploy the specified Helm chart using the values from your public Gist.

## Uninstall

To remove the Helmless Operator and its custom resources:

1. Delete the HelmLess custom resource:

```sh
kubectl delete -f helmless-deployer.yaml
```

2. Uninstall the operator's deployment and related resources:

```sh
kubectl delete -f helmless-operator.yaml
```

3. Delete the Helmless CRD:

```sh
kubectl delete -f config/crd/bases/cache.redrvm.cloud_helmlesses.yaml
```

This will clean up all resources related to the Helmless Operator.

## License

Copyright 2023.

Licensed under the MIT License (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://opensource.org/licenses/MIT

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.