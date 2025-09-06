package main

import (
	"fmt"

	"coldheater/internal/config"
	"coldheater/internal/database"
	"coldheater/internal/ui/cli"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Couldn't connect to a database:\n %v\n", err)
		return
	}
	defer db.Close()

	config.Global, err = config.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
	}

	fmt.Printf("Successfully connected to database\n")

	cli.ShowMenu(db)
}

