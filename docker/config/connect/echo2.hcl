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
  connect = {
    sidecar_service = {
      address = "10.3.30.12"
      port = 20000

      checks = [
        {
          name = "Connect Sidecar Listening"
          tcp = "10.3.30.12:20000"
          Interval = "10s"
        }
      ]

      proxy = {
        destination_service_name = "echo2-http"
        destination_service_id = "echo2-http"
        local_service_address = "10.3.0.12"
        local_service_port = 8080
      }
    }
  }
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
