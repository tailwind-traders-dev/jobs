# Deploy jobs using Kubernetes Event-driven Autoscaling (KEDA) on Azure Kubernetes Service (AKS)

After deployment, the below steps will walk you through the deployment of `jobs` to Azure Kubernetes Service (AKS) and the autoscaling of those jobs using Kubernetes Event-driven Autoscaling (KEDA).

## Install kubectl

Using the Azure CLI:

```bash
sudo az aks install-cli
```

Alternatively you can install kubectl manually per the Kubernetes ducomentation: <https://kubernetes.io/docs/tasks/tools/install-kubectl-macos/>

Example for macOS on Apple Silicon:

```bash
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/arm64/kubectl"

sudo mv kubectl /usr/local/bin/kubectl
```

## Authenticate to AKS

```bash
az aks get-credentials \
    --resource-group 231000-aks \
    --name aks1 \
    --overwrite-existing
```

## Install KEDA

```bash
kubectl apply -f https://github.com/kedacore/keda/releases/download/v2.10.1/keda-2.10.1.yaml
```

See: <https://keda.sh/docs/2.10/deploy/#yaml>

## Create the Kubernetes Secret

```bash
# set environment variables so AZURE_SERVICEBUS_CONNECTION_STRING is populated
source _/env.sh

# create the kubernetes secret
kubectl create secret generic connection-string-secret --from-literal=connection-string=${AZURE_SERVICEBUS_CONNECTION_STRING}

# change to the azure-kubernetes-service directory
cd deploy/azure-kubernetes-service/

# deploy the manifest
kubectl apply -f jobs-hello.yaml
```

## Confirm the Deployment and KEDA ScaledObject are working

```bash
# deployment
kubectl describe deploy/jobs-hello

# scaledobject
kubectl describe scaledobject/jobs-hello-scaledobject
```

## Add items to the Service Bus Queue via the Azure Portal

Use the Service Bus Explorer in the Azure Portal portal to send 100 "hello, world" messages.

## Watch the Deployment scale up

```
$ kubectl get deploy
NAME         READY   UP-TO-DATE   AVAILABLE   AGE
jobs-hello   4/5     5            4           18m
```

## Confirm the output on one of the Deployment containers

```
kubectl logs deploy/jobs-hello
```

## Confirm resources requests and limits

You will notice that only 4 of the 5 pods have successfully deployed. Let's find out why.

Get the node name:

```
$ kubectl get nodes
NAME                                STATUS   ROLES   AGE     VERSION
aks-nodepool1-31930126-vmss000000   Ready    agent   4h24m   v1.27.3
```

Run `kubectl describe node` to see the allocated resources:

```
kubectl describe node aks-nodepool1-31930126-vmss000000
```

Output:

```
...
Allocated resources:
  (Total limits may be over 100 percent, i.e., overcommitted.)
  Resource           Requests      Limits
  --------           --------      ------
  cpu                1370m (72%)   11990m (631%)
  memory             4976Mi (97%)  9020Mi (177%)
  ephemeral-storage  0 (0%)        0 (0%)
  hugepages-1Gi      0 (0%)        0 (0%)
  hugepages-2Mi      0 (0%)        0 (0%)
...
```

You can see `Requests` does not have enough CPU to schedule the 5th pod which has a CPU request of 100m.

Our next step will be to adjust the `requests` on our `Deployment` and/or the `maxReplicaCount` on our `ScaledObject` so that we can efficiently us the resources in our cluster.

In the meantime, let's wait for the 4 pods to complete their work.

## Confirm Deployment has scaled to zero

Run `kubectl get deploy`

```
kubectl get deploy
```

Output:

```
NAME         READY   UP-TO-DATE   AVAILABLE   AGE
jobs-hello   0/0     0            0           25m
```

The jobs have completed, the Deployment has scaled to zero, and this means the Azure Service Bus Queue is empty.
