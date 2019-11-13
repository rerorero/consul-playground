kind = "proxy-defaults"
name = "global"

config {
  # (dog)statsd listener on either UDP or Unix socket. 
  # envoy_statsd_url = "udp://127.0.0.1:9125"
  # envoy_dogstatsd_url = "udp://127.0.0.1:9125"

  # IP:port to expose the /metrics endpoint on for scraping.
  prometheus_bind_addr = "0.0.0.0:9102"

  # The flush interval in seconds.
  envoy_stats_flush_interval = 10
}
