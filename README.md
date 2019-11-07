consul-playground
==========

#### Docker Compose

plain configuration
```
docker-compose -f ./docker/docker-compose.yaml up -d
curl localhost:8000 -d 'alice'
```

service mesh with Consul proxy
```
docker-compose -f ./docker/docker-compose.connect.yaml up -d
curl localhost:8000 -d 'alice'
```


#### kubernetes
plain configuration
```
kubectl apply -f kube/plain.yaml
curl http://35.236.184.250 -d 'alice'
```

