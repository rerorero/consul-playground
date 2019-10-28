package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

// Env is env
type Env struct {
	ID   string `envconfig:"ID"`
	Port int    `envconfig:"PORT"`
}

func (s *Env) echo(w http.ResponseWriter, r *http.Request) {
	log.Printf("[echo:%s] receive request\n", s.ID)
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

func main() {
	var env Env
	envconfig.Process("", &env)

	http.HandleFunc("/", env.echo) // hello
	http.ListenAndServe(fmt.Sprintf(":%d", env.Port), nil)
}
