# edgex-on-kubernetes - WIP
Scripts to deploy EdgeX on Kubernetes (https://kubernetes.io). Work in progress. Not ready for deployment.

## Steps

### 1. Setup Consul

```
kubectl apply -f k8s/services/consul-service.yaml

kubectl apply -f k8s/deployments/consul-deployment.yaml
```

### 2. Create Config-Seed job

```
kubectl create -f k8s/jobs/config-seed-job.yaml
```

### 3. Setup MongoDB

```
kubectl apply -f k8s/services/mongo-service.yaml

kubectl apply -f k8s/deployments/mongo-deployment.yaml

```

### 4. Septup Mosquito MQTT Broker

Only if needed (external broker could be used instead)

```
kubectl apply -f k8s/services/mosquitto-service.yaml

kubectl apply -f k8s/deployments/mosquitto-deployment.yaml

```


### 5. Setup Logging Service

```
kubectl apply -f k8s/services/logging-service.yaml

kubectl apply -f k8s/deployments/logging-deployment.yaml

```
### 6. Setup Notifications Service

```
kubectl apply -f k8s/services/notifications-service.yaml

kubectl apply -f k8s/deployments/notifications-deployment.yaml
```

### 7. Setup Core Metadata Service

```
kubectl apply -f k8s/services/metadata-service.yaml

kubectl apply -f k8s/deployments/metadata-deployment.yaml

```

### 8. Setup Core Data Service

```
kubectl apply -f k8s/services/data-service.yaml

kubectl apply -f k8s/deployments/data-deployment.yaml

```

### 9. Setup Core Command Service

```
kubectl apply -f k8s/services/command-service.yaml

kubectl apply -f k8s/deployments/command-deployment.yaml

```

### 10. Setup Scheduler Service

```
kubectl apply -f k8s/services/scheduler-service.yaml

kubectl apply -f k8s/deployments/scheduler-deployment.yaml

```

### 11. Setup Export Client Registration Service

```
kubectl apply -f k8s/services/export-client-service.yaml

kubectl apply -f k8s/deployments/export-client-deployment.yaml

```

### 12. Setup Export Distribution Service 

```
kubectl apply -f k8s/services/export-distro-service.yaml

kubectl apply -f k8s/deployments/export-distro-deployment.yaml

```
### 13. Setup Rules Engine Service

```
kubectl apply -f k8s/services/rules-engine-service.yaml

kubectl apply -f k8s/deployments/rules-engine-deployment.yaml

```
### 14. Setup Rules Engine Service

```
kubectl apply -f k8s/services/rules-engine-service.yaml

kubectl apply -f k8s/deployments/rules-engine-deployment.yaml

```

### 15. Create MetalLB  L2 Load Balancer to provide external access to Mainflux Services

For more information see [MetalLB L2 tutorial](https://metallb.universe.tf/tutorial/layer2/)

```
kubectl apply -f k8s/metallb/metallb.yaml

kubectl apply -f k8s/metallb/layer2-config.yaml
```

### 16. Configure Internet access
Configure NAT on your Firewall to forward port 1883 (MQTT)

