service {
  name = "echo-proxy"
  tags = ["echo-proxy","http"]
  address = "10.3.0.13"
  port = 8000
  checks = [{
    tcp = "10.3.0.13:8000"
    interval = "10s"
    timeout = "5s"
  }]
}
