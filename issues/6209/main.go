package main

import (
	"fmt"

	"github.com/hashicorp/consul/agent/config"
)

func main() {
	test1 := `
{
  "data_dir": "./tmpdata"
  "ports": {
    "dns": 18600
  },
}
`
	c, err := config.Parse(test1, "json")

	fmt.Printf("test1 error: %v\n", err)
	fmt.Printf("%#v", c)
}
