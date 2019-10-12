<p align="center"><img src="./img/kafka_logo.png" width="340"></p>

# Transport Nginx Access Logs into Kafka with Logging Operator

<p align="center"><img src="./img/nignx-kafka.png" width="900"></p>

---
## Contents
- **Installation**
  - Kafka 
    - [Deploy with Helm](#deploy-kafka)
  - **Logging Operator**
    - [Deploy with Helm](#install-with-helm)
    - [Deploy with Kuberenetes Manifests](./deploy/README.md#deploy-logging-operator-from-kubernetes-manifests)
   - **Demo Application**  
    - [Deploy with Helm](install-with-helm)
    - [Deploy with Kuberenetes Manifests](#install-from-kubernetes-manifests)
- **Validation**
    - [Kafkacat](#test-your-deployment-with-kafkacat)
---

## Deploy Kafka
>In this demo we are using our kafka operator.
> [Easy Way Installing with Helm](https://github.com/banzaicloud/kafka-operator#easy-way-installing-with-helm)
## Deploy Logging-Operator with Demo Application

### Install with Helm 

[Install Logging-operator with helm](./deploy/README.md#deploy-logging-operator-with-helm)


#### Nginx App and Logging Definition
```bash
helm install --namespace logging --name nginx-demo banzaicloud-stable/nginx-logging-kafka-demo
```

---
### Install from Kubernetes manifests
[Install Logging-operator from manifests](./deploy/README.md#deploy-logging-operator-from-kubernetes-manifests)

#### Create `logging` resource
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: Logging
metadata:
  name: default-logging-simple
spec:
  fluentd: {}
  fluentbit: {}
  controlNamespace: logging
EOF
```

> Note: `ClusterOutput` and `ClusterFlow` resource will only be accepted in the `controlNamespace` 


#### Create an Kafka output definition 
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: Output
metadata:
  name: kafka-output
spec:
  kafka:
    brokers: kafka-headless.kafka.svc.cluster.local:29092
    default_topic: topic
    format: 
      type: json    
    buffer:
      tags: topic
      timekey: 1m
      timekey_wait: 30s
      timekey_use_utc: true
EOF
```
> Note: For production set-up we recommend using longer `timekey` interval to avoid generating too many object.

#### Create `flow` resource
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: logging.banzaicloud.io/v1beta1
kind: Flow
metadata:
  name: kafka-flow
spec:
  filters:
    - parser:
        key_name: message
        remove_key_name_field: true
        reserve_data: true
        parsers:
          - type: nginx
  selectors:
    app: nginx
  outputRefs:
    - kafka-output
EOF
```

#### Install nginx deployment
```bash
cat <<EOF | kubectl -n logging apply -f -
apiVersion: apps/v1 
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
          name: http
          protocol: TCP
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
EOF
```

### Test Your Deployment with kafkacat
##### Exec Kafaka test pod
```bash
kubectl -n kafka exec -it kafka-test-c sh
```
Run kafkacat
```bash
kafkacat -C -b kafka-0.kafka-headless.kafka.svc.cluster.local:29092 -t topic
```

[![asciicast](https://asciinema.org/a/273236.svg)](https://asciinema.org/a/273236)
