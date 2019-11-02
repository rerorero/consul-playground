services {
  name = "echo2-http"
  tags = ["echo","http"]
  address = "10.3.0.12"
  port = 8080
  checks = [{
    tcp = "10.3.0.12:8080"
    interval = "10s"
    timeout = "5s"
  }]
}

services {
  name = "echo2-grpc"
  tags = ["echo","grpc"]
  address = "10.3.0.12"
  port = 9090
  checks = [{
    tcp = "10.3.0.12:9090"
    interval = "10s"
    timeout = "5s"
  }]
}
