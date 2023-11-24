package main

import (
	"flag"
	"fmt"
)

func main() {
	port := flag.String("port", "8080", "port to serve on")
	envID := flag.String("env-id", "", "environment to route connections")
	defaultAddr := flag.String("default-addr", "", "default address")

	flag.Parse()

	fmt.Println("port:", *port)
	fmt.Println("env-id:", *envID)
	fmt.Println("default-addr:", *defaultAddr)
}
