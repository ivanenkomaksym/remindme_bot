package main

import (
	"github.com/ivanenkomaksym/remindme_bot/api/route"
	"github.com/ivanenkomaksym/remindme_bot/bootstrap"
)

func main() {
	// Initialize the application
	app := bootstrap.App()

	// Setup routes and start the server
	route.Setup(&app)
}
