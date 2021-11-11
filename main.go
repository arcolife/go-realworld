package main

import (
	"log"

	"github.com/0xdod/go-realworld/postgres"
	"github.com/0xdod/go-realworld/server"
)

func main() {
	port := ":8000"
	db, err := postgres.Open("postgres://admin:admin@localhost:5432/conduit?sslmode=disable")
	if err != nil {
		log.Fatalf("cannot open database: %v", err)
	}

	srv := server.NewServer(db)
	log.Fatal(srv.Run(port))
}
