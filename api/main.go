package main

import (
	"Codex-Backend/api/aws"
	"Codex-Backend/api/server"
	"log"
)

func main() {

	svc, err := aws.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	server.Server(svc)
}
