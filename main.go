package main

import (
	"fmt"
	"log"

	"github.com/vla47/course-api/router"
)

func main() {

	// setup the server and start the listener
	server := router.LoadRoutes()

	log.Fatal(server)
	fmt.Scanln() 
}
