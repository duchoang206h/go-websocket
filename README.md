# Golang websocket sample with k8s, hpa

## Deployment
*** Kubectl:
```bash
kubectl create namespace websocket
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml # metrics-server
kubectl apply -f ./deployment/redis.yaml -n websocket #  redis
kubectl apply -f ./deployment/websocket.yaml -n websocket # go-websocket
kubectl apply -f ./deployment/hap.yaml -n websocket # hpa
go run client.go # spam traffic to go-websocket
```

