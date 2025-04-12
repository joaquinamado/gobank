package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/joaquinamado/gobank/internal/app/api"
	. "github.com/joaquinamado/gobank/internal/app/storage"
	. "github.com/joaquinamado/gobank/internal/app/types"
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

//	@title			GoBank API
//	@version		1.0
//	@description	An API for a simple bank
//	@termsOfService	None

//	@contact.name	API Support
//	@contact.url	None
//	@contact.email	None

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/v1

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	seed := flag.Bool("seed", false, "Seed the database")
	flag.Parse()

	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding database")
		// Seed stuff
		seedAccounts(store)
	}

	server := NewApiServer(":3000", store)
	mux := server.Mount()
	server.Run(mux)
}
