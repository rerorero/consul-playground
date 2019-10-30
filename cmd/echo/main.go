package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	echo "github.com/rerorero/consul-playground/proto"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// Env is env
type Env struct {
	ID                 string   `envconfig:"ID"`
	HTTPPort           *int     `envconfig:"HTTP_PORT"`
	GRPCPort           *int     `envconfig:"GRPC_PORT"`
	ProxyHTTPPort      *int     `envconfig:"PROXY_HTTP_PORT"`
	ProxyHTTPUpstreams []string `envconfig:"HTTP_UPSTREAMS"`
	ProxyGRPCUpstreams []string `envconfig:"GRPC_UPSTREAMS"`
	Insecure           bool     `envconfig:"INSECURE" default:"true"`
}

// echo is HTTP echo server handler
func (s *Env) echo(w http.ResponseWriter, r *http.Request) {
	log.Printf("[echo:%s] receive HTTP request\n", s.ID)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprintf(w, "Hi, %s. I am %s HTTP", string(body), s.ID); err != nil {
		log.Printf("[echo:%s] Error %v\n", s.ID, err)
	}
}

// Echo is GRPC echo server handler
func (s *Env) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	log.Printf("[echo:%s] receive GRPC request\n", s.ID)
	return &echo.EchoResponse{
		Message: fmt.Sprintf("Hi, %s. I am %s GRPC", req.Message, s.ID),
	}, nil
}

// ProxyServer is HTTP proxy server handler
type ProxyServer struct {
	httpClient  *http.Client
	gRPCClients []echo.EchoClient
	env         Env
}

func newProxy(env Env) *ProxyServer {
	httpClient := new(http.Client)

	var grpcClients []echo.EchoClient
	var opts []grpc.DialOption
	if env.Insecure {
		opts = append(opts, grpc.WithInsecure())
	}
	for _, upstream := range env.ProxyGRPCUpstreams {
		conn, err := grpc.Dial(upstream, opts...)
		if err != nil {
			log.Fatalf("did not connect to %s %v", upstream, err)
		}
		grpcClients = append(grpcClients, echo.NewEchoClient(conn))
	}

	return &ProxyServer{
		httpClient:  httpClient,
		gRPCClients: grpcClients,
		env:         env,
	}
}

func (s *ProxyServer) proxy(w http.ResponseWriter, r *http.Request) {
	log.Printf("[proxy:%s] receive request\n", s.env.ID)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var responseBody string

	// choice randomly
	random := rand.Intn(len(s.env.ProxyHTTPUpstreams) + len(s.gRPCClients))
	if random < len(s.env.ProxyHTTPUpstreams) {
		// HTTP
		host := s.env.ProxyHTTPUpstreams[random]
		log.Printf("[proxy:%s] send HTTP request to %s\n", s.env.ID, host)
		res, err := s.httpClient.Post(host, "text/plain", bytes.NewReader(body))
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if res.StatusCode != http.StatusOK {
			log.Printf("[proxy:%s] HTTP responds error=%v\n", s.env.ID, res.StatusCode)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		responseBody = string(b)

	} else {
		// gRPC
		log.Printf("[proxy:%s] send gRPC request to %s\n", s.env.ID, s.env.ProxyGRPCUpstreams[random-len(s.env.ProxyHTTPUpstreams)])
		cli := s.gRPCClients[random-len(s.env.ProxyHTTPUpstreams)]
		res, err := cli.Echo(r.Context(), &echo.EchoRequest{
			Message: string(body),
		})
		if err != nil {
			log.Printf("[proxy:%s] gRPC responds error=%v\n", s.env.ID, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseBody = res.Message
	}

	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, responseBody); err != nil {
		log.Printf("[echo:%s] Error %v\n", s.env.ID, err)
	}
}

func main() {
	rand.Seed(time.Now().Unix())
	var env Env
	envconfig.Process("", &env)

	log.Printf("[echo:%s] started", env.ID)

	var wg errgroup.Group

	// HTTP
	if env.HTTPPort != nil {
		wg.Go(func() error {
			http.HandleFunc("/", env.echo) // hello
			return http.ListenAndServe(fmt.Sprintf(":%d", *env.HTTPPort), nil)
		})
	}

	// GRPC
	if env.GRPCPort != nil {
		wg.Go(func() error {
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *env.GRPCPort))
			if err != nil {
				log.Fatalf("[echo:%s] failed to listen: %v", env.ID, err)
			}
			s := grpc.NewServer()
			echo.RegisterEchoServer(s, &env)
			return s.Serve(lis)
		})
	}

	// proxy
	if env.ProxyHTTPPort != nil {
		wg.Go(func() error {
			proxyServer := newProxy(env)
			http.HandleFunc("/", proxyServer.proxy)
			return http.ListenAndServe(fmt.Sprintf(":%d", *env.ProxyHTTPPort), nil)
		})
	}

	if err := wg.Wait(); err != nil {
		log.Printf("error occured: %v\n", err)
	}
}
