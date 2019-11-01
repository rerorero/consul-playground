consul-playground
==========

#### Docker Compose

plain configuration
```
docker-compose -f ./docker/docker-compose.yaml up -d
curl localhost:8000 -d 'alice'
```

#### kubernetes
plain configuration
```
kubectx <your context>
kubens <your namespace>
kubectl apply -f kube/plain.yaml
```

