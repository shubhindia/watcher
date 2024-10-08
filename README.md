![Watcher](img/the-watcher-marvel.gif)

# This one *INTERVENES*

## Motivation
This project was created to help me felicitate eureka-server restarts whenever underlying services are restarted. There is a need to restart eureka-server whenever
a service is restarted to ensure that the eureka-server is aware of the changes in the service registry. Setting up pod priorities helped in initial starts but it did
not help with random pod restarts. This project is a simple watcher that watches for pods and restarts the eureka-server whenever a pod is restarted to make sure that
the eureka-server is aware of the changes in the service registry.

## How it works
The watcher is a simple kubernetes controller written in golang. It watches pod creation events in given namespace and restarts the desired pods when a new pod is created.

## How to use (WIP)
1. Create the config 
```yaml
namespace: <namespace>
newest:
  - Kind: <Deployment/StatefulSet>
    Name: <Name>
```
2. Mount it as a configmap
3. Deploy the watcher