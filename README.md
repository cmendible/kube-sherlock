[![Build Status](https://mendible.visualstudio.com/mendible/_apis/build/status/cmendible.kube-sherlock?branchName=master)](https://mendible.visualstudio.com/mendible/_build/latest?definitionId=8&branchName=master)

# kube-sherlock

kube-sherlock lists all pods which do not have the labels listed in the **config.yaml** file.

The default **config.yaml** values are:

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

## Sample results

``` shell
+------------------------------+-------------+-----------------------------------------------------------------+
|            LABEL             |  NAMESPACE  |                            POD NAME                             |
+------------------------------+-------------+-----------------------------------------------------------------+
| app.kubernetes.io/version    | default     | mypod                                                           |
+                              +-------------+-----------------------------------------------------------------+
|                              | kube-system | aci-connector-linux-79b768b6d6-fhb9d                            |
+                              +             +-----------------------------------------------------------------+
|                              |             | addon-http-application-routing-default-http-backend-5ccb95j9dgb |
```
