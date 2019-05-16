# Onetime Password Adapter

> This process is based off:
> - https://github.com/salrashid123/istio_custom_auth_adapter
> - https://github.com/istio/istio/wiki/Mixer-Out-Of-Process-Adapter-Walkthrough

## Introduction
This is an out of process Mixer Adapter that handles authorization checks. In this case the adapter depends on another API Builder service. In the future the actual adapter could be the API Builder process.

## Building
The `build.sh` runs through the steps of building everything and running a quick test look at it for info.

## Deploying configurations
```
kubectl apply -f $ROOT_FOLDER/generated/attributes.yaml -f $ROOT_FOLDER/generated/template.yaml
kubectl apply -f $ROOT_FOLDER/generated/onetimeadapter.yaml
```

## Deploying adapter
### Publish Image
```
docker build -t <your dockerhub>/onetimeadapter .
docker push <your dockerhub>/onetimeadapter
```

### Install Adapter
```
kubectl apply -f $ROOT_FOLDER/onetime-service.yaml
```

## Testing
Install a rule that uses the adapter as a handler. This example applies the rule to all inbound requests for services in the `apic-demo` namespace.
```
match: source.labels["istio"] == "ingressgateway" && destination.namespace == "apic-demo"
```

It expects two headers - `user` identifies the user and `apikey` is the token to verify.

```
 params:
   subject:
     user: request.headers["user"]
     properties:
       custom_token_header: request.headers["apikey"]
```

To install:

```
kubectl apply -f $ROOT_FOLDER/onetime-apic-demo.yaml
```

