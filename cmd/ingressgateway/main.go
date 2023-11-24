package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/picop-rd/ingressgateway/app/ingressgateway"
)

func main() {
	port := flag.String("port", "8000", "port to serve on")
	envID := flag.String("env-id", "", "environment to route connections")
	destination := flag.String("destination", "", "destination address")

	flag.Parse()

	server := ingressgateway.New(*envID, *destination)
	go server.Start(fmt.Sprintf(":%s", *port))
	defer server.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
