package main

import (
	"log"

	"github.com/gurrpi/kind-manager/server"
)

func main() {
	s := server.New()
	if err := s.Run(); err != nil {
		log.Fatal("Server couldn't started: ", err)
	}
}
