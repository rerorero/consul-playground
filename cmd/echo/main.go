package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"

	echo "github.com/rerorero/consul-playground/proto"

	"google.golang.org/grpc"

	"github.com/kelseyhightower/envconfig"
)

// Env is env
type Env struct {
	ID       string `envconfig:"ID"`
	HTTPPort int    `envconfig:"HTTP_PORT"`
	GRPCPort int    `envconfig:"GRPC_PORT"`
}

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
	fmt.Fprintf(w, "%s from %s", string(body), s.ID)
}

func (s *Env) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	log.Printf("[echo:%s] receive GRPC request\n", s.ID)
	return &echo.EchoResponse{
		Message: fmt.Sprintf("%s from %s", req.Message, s.ID),
	}, nil
}

func main() {
	var env Env
	envconfig.Process("", &env)

	log.Printf("[echo:%s] started", env.ID)

	var wg errgroup.Group

	// HTTP
	wg.Go(func() error {
		http.HandleFunc("/", env.echo) // hello
		return http.ListenAndServe(fmt.Sprintf(":%d", env.HTTPPort), nil)
	})

	// GRPC
	wg.Go(func() error {

		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", env.GRPCPort))
		if err != nil {
			log.Fatalf("[echo:%s] failed to listen: %v", env.ID, err)
		}
		s := grpc.NewServer()
		echo.RegisterEchoServer(s, &env)
		return s.Serve(lis)
	})

	if err := wg.Wait(); err != nil {
		log.Printf("error occured: %v\n", err)
	}
}
