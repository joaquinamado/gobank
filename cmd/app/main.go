package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"

	. "github.com/joaquinamado/gobank/internal/app/api"
	"github.com/joaquinamado/gobank/internal/app/env"
	"github.com/joaquinamado/gobank/internal/app/repositories"
	"github.com/joaquinamado/gobank/internal/app/services"
)

func Init() {
	// Set main function to run on the main thread.
	runtime.LockOSThread()
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
	go startServer()

	// Listen for functions that need to run on the main thread.
	var quit = make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	for {
		select {
		case f := <-services.Run:
			f()
		case <-quit:
			log.Println("shutting down")
			return
		}
	}
}

func startServer() {
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
