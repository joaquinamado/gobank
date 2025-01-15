package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/joaquinamado/gobank/internal/app/api"
	. "github.com/joaquinamado/gobank/internal/app/storage"
	. "github.com/joaquinamado/gobank/internal/app/types"
	. "github.com/joaquinamado/gobank/internal/app/utils"
)

func seedAccount(store Storage, fname, lname, pw string) *Account {
	acc, err := NewAccount(fname, lname, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("new account created => ", acc.Number)

	return acc
}

func seedAccounts(store Storage) {
	seedAccount(store, "John", "Doe", "password")
}

func main() {
	seed := flag.Bool("seed", false, "Seed the database")
	flag.Parse()

	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	// Creates the singleton instance of the .env file
	GetEnvInstance()

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding database")
		// Seed stuff
		seedAccounts(store)
	}

	server := NewApiServer(":3000", store)
	server.Run()
}
