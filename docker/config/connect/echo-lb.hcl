service {
  name = "echo-lb"
  tags = ["echo-lb","http"]
  address = "10.3.0.13"
  port = 8000
  checks = [{
    tcp = "10.3.0.13:8000"
    interval = "10s"
    timeout = "5s"
  }]

  connect = {
    sidecar_service = {
      address = "10.3.30.13"
      port = 20000

      checks = [
        {
          name = "Connect Sidecar Listening"
          tcp = "10.3.30.13:20000"
          Interval = "10s"
        }
      ]

      proxy = {
        destination_service_name = "echo-lb"
        destination_service_id = "echo-lb"
        local_service_address = "10.3.0.13"
        local_service_port = 8000
        upstreams = [
          {
            destination_type = "service"
            destination_name = "echo1-http"
            local_bind_address = "10.3.30.13"
            local_bind_port = 8001
          }
        ]
      }
    }
  }
}
