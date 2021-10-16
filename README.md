[![kube-sherlock](https://github.com/cmendible/kube-sherlock/actions/workflows/build.yaml/badge.svg)](https://github.com/cmendible/kube-sherlock/actions/workflows/build.yaml)

# kube-sherlock

## Check Pod Requests and Limits

To check pod requests and limits in the `kube-system` namespaces, you can use the following command:

``` shell
./kube-sherlock resources -n kube-system -kubeconfig
```

## Check Missing Labels

To check which pods do not have the `app.kubernetes.io/name` label in the `kube-system` namespace:

``` shell
./kube-sherlock.sh labels -l app.kubernetes.io/name -n kube-system -kubeconfig
```

Using a config file:

kube-sherlock lists all pods which do not have the labels listed in the **config.yaml** file.

The default **config.yaml** values are:

``` yaml
labels:
  - "app.kubernetes.io/name"
  - "app.kubernetes.io/instance"
  - "app.kubernetes.io/version"
  - "app.kubernetes.io/component"
  - "app.kubernetes.io/part-of"
  - "app.kubernetes.io/managed-by"
```

It's also possible to specify the namespaces you want to scan in the **config.yaml**:

``` yaml
namespaces:
  - default
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
