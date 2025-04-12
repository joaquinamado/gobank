package main

import (
	"flag"
	"fmt"
	"log"

	. "github.com/joaquinamado/gobank/internal/app/api"
	"github.com/joaquinamado/gobank/internal/app/env"
	"github.com/joaquinamado/gobank/internal/app/repositories"
)

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

	repo, err := repositories.NewRepository()

	if err != nil {
		log.Fatal(err)
	}

	if *seed {
		fmt.Println("Seeding database")
		// Seed stuff
		repo.Account.SeedAccounts()
	}

	port := env.GetString("API_PORT", "8080")

	server := NewApiServer(port, *repo)
	mux := server.Mount()
	server.Run(mux)
}
