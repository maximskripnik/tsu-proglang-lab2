package main

import (
	"client"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"server"
)

func main() {

	parser := argparse.NewParser("", "Prints provided string to stdout")
	mode := parser.Selector("m", "mode", []string{"server", "client"}, &argparse.Options{Required: true, Help: "Mode"})
	port := parser.Int("p", "port", &argparse.Options{Required: true, Help: "For server mode - tcp port to lisen. For client - server port to connect to"})
	maxConnections := parser.Int("n", "max-connections", &argparse.Options{Required: false, Help: "Number of maximum allowed connections. Only makes sense for the server mode", Default: 100})
	host := parser.String("", "host", &argparse.Options{Required: false, Help: "Hostname of the server to connect to. Only makes sense for the client mode"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Print(parser.Usage(err))
	}

	switch *mode {
	case "server":
		server.Serve(*port, *maxConnections)
	case "client":
		client.Connect(*host, *port)
	}

}
