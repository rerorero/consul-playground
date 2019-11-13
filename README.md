consul-playground
==========

## Docker Compose

#### plain configuration
```
docker-compose -f ./docker/docker-compose.yaml up -d
curl localhost:8000 -d 'alice'
```

#### service mesh with Consul proxy
```
docker-compose -f ./docker/docker-compose.connect.yaml up -d
curl localhost:8000 -d 'alice'
```

#### service mesh with Envoy proxy
```
docker-compose -f ./docker/docker-compose.envoy.yaml up -d
curl localhost:8000 -d 'alice'
```


## kubernetes

#### plain configuration
```
kubectl apply -f kube/echo.plain.yaml
curl http://35.236.184.250 -d 'alice'
```

#### service mesh with Envyo proxy
```
helm init

# if RBAC is enabled
kubectl apply -f kube/admin-tiller.yaml
kube patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'

kubectl apply -f kube/echo.connect.yaml
curl http://35.236.184.250 -d 'alice'
```

#### Prometheus integration
```
helm install --name prometheus --namespace default -f kube/prometheus-values.yaml stable/prometheus
```

#### Show control plane UI
```
kubectl port-forward consul-playground-consul-server-0 18500:8500
open http://localhost:18500/ui/dc10/services
```

#### Confirm Envoy configuration
```
# you can replace the app lable
kubectl port-forward `kubectl get po -l app=echo-lb -o jsonpath="{.items[0].metadata.name}"` 19000:19000
open http://localhost:19000
```
