package sender

import (
	"coldheater/internal/database"
	"fmt"
)

func DebugMenu(){
	options := [...]string{"Cold email from a specific sender",}

	for i, option := range options{
		fmt.Printf("%d. %s\n", i+1, option)
	}

	var input int
	fmt.Scanln(&input)
	switch input{
	case 1:
		coldEmailRaw();
	default:
		fmt.Println("incorrect input")
	}
}

func coldEmailRaw(){

	var input string

	fmt.Println("Sender email:")
	fmt.Scanln(&input)
	testSender := database.Sender{Email: input}

	fmt.Println("Reciever email:")
	fmt.Scanln(&input)
	testReciever := input
	testContent := EmailContent{Body: "If you see this message, sending is working fine", Subject: "Hello world!"}

	err := SendColdEmail(testSender, testReciever, testContent)
	if err != nil {
		fmt.Println(err)
	}
}