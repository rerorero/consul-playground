services {
  id = "echo1-http"
  name = "echo-http"
  tags = ["echo","http"]
  address = "10.3.0.11"
  port = 8080
  checks = [{
    tcp = "10.3.0.11:8080"
    interval = "10s"
    timeout = "5s"
  }]
  connect = {
    sidecar_service = {
      address = "10.3.30.11"
      port = 20000

      checks = [
        {
          name = "Connect Sidecar Listening"
          tcp = "10.3.30.11:20000"
          Interval = "10s"
        }
      ]

      proxy = {
        destination_service_name = "echo-http"
        destination_service_id = "echo1-http"
        local_service_address = "10.3.0.11"
        local_service_port = 8080
      }
    }
  }
}

services {
  id = "echo1-grpc"
  name = "echo-grpc"
  tags = ["echo","grpc"]
  address = "10.3.0.11"
  port = 9090
  checks = [{
    tcp = "10.3.0.11:9090"
    interval = "10s"
    timeout = "5s"
  }]

  connect = {
    sidecar_service = {
      address = "10.3.40.11"
      port = 20000

      checks = [
        {
          name = "Connect Sidecar Listening"
          tcp = "10.3.40.11:20000"
          Interval = "10s"
        }
      ]

      proxy = {
        destination_service_name = "echo-grpc"
        destination_service_id = "echo1-grpc"
        local_service_address = "10.3.0.11"
        local_service_port = 9090
      }
    }
  }
}
