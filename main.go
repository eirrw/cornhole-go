package main

import (
	"log"
	"os"
	"virunus.com/cornhole/config"
	"virunus.com/cornhole/database"
	"virunus.com/cornhole/server"
	"virunus.com/cornhole/tui"
)

func main() {
	configuration, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "init":
			database.InitializeDatabse(configuration)
			break
		case "serve":
			server.Serve(configuration)
			break
		}
	} else {
		app := tui.Get(configuration)
		if err = app.Run(); err != nil {
			log.Fatal(err)
		}
	}
}
