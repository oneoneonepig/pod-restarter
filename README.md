# pod-restarter
Restarts some Pods in Kubernetes, selected by labels

## Usage
```
-grace-period int
      the duration in seconds before the object should be deleted. (default 30)
-namespace string
      specify the namespace of the pods (default "default")
-selector string
      label selector
```

## Example
Delete all Pods in default Namespace
```
pod-restarter -namespace=default
```
Delete all Pods with label "app=nginx" in default Namespace
```
pod-restarter -namespace=default -selector="app=nginx"
```

## Library and reference
- https://github.com/kubernetes/client-go
- https://github.com/kubernetes/client-go/blob/master/examples/out-of-cluster-client-configuration/main.go
- https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go

