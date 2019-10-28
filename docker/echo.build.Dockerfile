FROM golang:1.13-stretch
WORKDIR /go/src/github.com/rerorero/consul-playground
COPY . ./
RUN make build

FROM alpine:latest
COPY --from=0 /go/src/github.com/rerorero/consul-playground/bin/echo /bin/echo
CMD ["/bin/echo"]
