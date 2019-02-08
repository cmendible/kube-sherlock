[![Build Status](https://mendible.visualstudio.com/mendible/_apis/build/status/cmendible.kube-sherlock?branchName=master)](https://mendible.visualstudio.com/mendible/_build/latest?definitionId=8&branchName=master)

# kube-sherlock

kube-sherlock lists all pods which do not have the lables listed in the **config.yaml** file.

config.yaml default values are:

``` shell
labels:
  - "app.kubernetes.io/name"
  - "app.kubernetes.io/instance"
  - "app.kubernetes.io/version"
  - "app.kubernetes.io/component"
  - "app.kubernetes.io/part-of"
  - "app.kubernetes.io/managed-by"
```

## Running in a Kubernetes cluster without RBAC enabled

``` shell
kubectl run --rm -i -t kube-sherlock --image=cmendibl3/kube-sherlock:0.1 --restart=Never
```

## Running in a Kubernetes cluster with RBAC enabled

``` shell
kubectl apply -f service-account.yaml
kubectl run --rm -i -t kube-sherlock --image=cmendibl3/kube-sherlock:0.1 --restart=Never --overrides='{ \"apiVersion\": \"v1\", \"spec\": { \"serviceAccountName\": \"kube-sherlock\" } }'
```

## Sample results:

``` shell
+------------------------------+-------------+-----------------------------------------------------------------+
| app.kubernetes.io/instance   | default     | kube-sherlock                                                   |
+                              +             +-----------------------------------------------------------------+
|                              |             | mypod    
```