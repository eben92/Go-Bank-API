package main

import (
	"flag"
	"fmt"
	"log"
)

func seedAccount(s Storage, fname, lname, password string) *Account {
	acc, err := NewAccount(fname, lname, password)

	if err != nil {
		log.Fatal(err)
	}

	if err := s.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account =>", acc.Number)

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "Daniel", "Opsman", "12333")
}

func main() {

	seed := flag.Bool("seed", false, "Seed the db")
	flag.Parse()

	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {

		// cmd: "./bin/gobank --seed"

		fmt.Println("seeding the database")
		seedAccounts(store)
		// store.DropDB()
	}

	server := NewAPIServer(":3000", store)
	server.Run()
}

// 8449
