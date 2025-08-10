package cli

import "fmt"

func ShowMenu() {
	options := [...]string{"Create new bot", "Upload leads", "Start cold outreach", "Exit"}

	for {
		fmt.Printf("=== ColdHeater CLI ===\n")
		for i, name := range options {
			fmt.Printf("%d. %s\n", i+1, name)
		}

		var userInput uint8

		fmt.Printf("Enter digit to select an option (e.g. \"%d\")\n", len(options))
		n, err := fmt.Scanf("%d", &userInput)
		if n != 1 || err != nil {
			fmt.Printf("Wasnt able to read user's input. Example usage: \"4\"[Enter]:\n %v\n", err)
		}

		switch userInput {
		case 1:
			fmt.Printf("Creating new bot...\n")
		case 2:
			fmt.Printf("Uploading leads...\n")
		case 3:
			fmt.Printf("Starting cold outreach...\n")
		case 4:
			fmt.Printf("Dust to dust...\n")
			return
		default:
			fmt.Printf("Invalid input. Enter a digit between 1-%d\n", len(options))
		}
	}
}
