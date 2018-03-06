package main

import (
  "internal/api"
  "flag"
)

func main() {
  port := flag.String("port", "8000", "exposed port of blockchain node.")
  flag.Parse()
  api.StartNode(port)
}
