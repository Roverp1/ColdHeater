package cli

import (
	botcreation "coldheater/internal/bot_creation"
	"coldheater/internal/database"
	"coldheater/internal/sender"
	"database/sql"
	"fmt"
	"os"
)

func ShowMenu(db *sql.DB) {
	options := [...]string{"Create new bot", "Upload leads", "Start cold outreach", "List all bots", "DDbug", "Exit"}

	for {
		fmt.Printf("=== ColdHeater CLI ===\n")
		for i, name := range options {
			fmt.Printf("%d. %s\n", i+1, name)
		}

		var userInput uint8

		fmt.Printf("Enter digit to select an option (e.g. \"%d\")\n", len(options))
		n, err := fmt.Scanln(&userInput)
		if n != 1 || err != nil {
			fmt.Printf("Wasnt able to read user's input. Example usage: \"4\"[Enter]:\n %v\n", err)
		}

		switch userInput {
		case 1:
			verificationAcc, err := database.GetVerificationAccount(db)
			if err != nil {
				fmt.Printf("Failed to get verification account:\n%v\n", err)
				return
			}
			bot, wasUsed, err := botcreation.CreateGmailBot(*verificationAcc)
			if err != nil {
				fmt.Printf("Failed to create bot:\n%v\n", err)
				return
			}
			if bot == nil {
				fmt.Printf("Failed to create bot")
				return
			}

			database.InsertBot(db, bot)

			if wasUsed == false {
				fmt.Printf("Verification account was not used for bot: %s", bot.Email)
				return
			}

			err = database.IncrementVerificationAccUsage(db, verificationAcc.Email)
			if err != nil {
				fmt.Printf("Failed to update verification account:\n%v\n", err)
				return
			}
		case 2:
			fmt.Printf("Uploading leads...\n")
		case 3:
			fmt.Printf("Starting cold outreach...\n")
		case 4:
			bots, err := database.GetAllBots(db)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			for _, bot := range bots {
				fmt.Println(bot)
			}
		case 5:
			sender.DebugMenu()
		case 6:
			fmt.Printf("Dust to dust...\n")
			return
		default:
			fmt.Printf("Invalid input. Enter a digit between 1-%d\n", len(options))
		}

		fmt.Printf("\nPress ENTER to continue...\n")
		var temp [1]byte
		os.Stdin.Read(temp[:])
	}
}
