package main

import (
	"fmt"

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

	fmt.Printf("Successfully connected to database\n")
}
