package main

import (
	"flag"
	"internal/api"
)

func main() {
	port := flag.String("port", "8000", "exposed port of blockchain node.")
	flag.Parse()
	api.StartNode(port)
}
