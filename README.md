# Kubernetes Admission Webhook

An example repo containing a [Kubernetes Amission Webhook](https://kubernetes.io/docs/admin/extensible-admission-controllers/#admission-webhooks).

This is based on the [upstream test](https://github.com/kubernetes/kubernetes/tree/release-1.10/test/images/webhook) with additional pieces required to build your own webhook. The tests can be referenced to created webhooks for other Kubernetes resources.

This projects contains:

- [x] A basic pod validator implementation that denis pods deployed with `hostNetwork=true`
- [x] Example TLS certificates
- [x] [Glide](https://github.com/Masterminds/glide) configuration to build the project with `client-go` and `apimachinery` dependencies
- [x] A Makefile to vendor and build the project inside a Docker container

## Prerequisites

- Ensure that the Kubernetes cluster is at least as new as v1.9.
- Ensure that `MutatingAdmissionWebhook` and `ValidatingAdmissionWebhook` admission controllers are enabled.
- Ensure that the `admissionregistration.k8s.io/v1beta1` API is enabled.



## Run It

#### Create TLS certs

``` sh
./scripts/pki.sh
```

#### Set ENV variables to be used when creating deployments

``` sh
export CA=`cat pki/example/ca.pem | base64`
export TLS_CERT=`cat pki/example/admission-webhook.pem | base64`
export TLS_KEY=`cat pki/example/admission-webhook-key.pem | base64`
```

#### Deploy Webhook

``` sh
./scripts/deploy.sh
```

#### Run example deployments

``` sh
kubectl apply -f examples/
```

#### Validate Admitter is checking pods

``` sh
➜  admission-webhook git:(master) ✗ kubectl get deploy
NAME                DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
admission-webhook   1         1         1            1           42s
nginx               1         1         1            1           26s
nginx-denied        1         0         0            0           26s
```