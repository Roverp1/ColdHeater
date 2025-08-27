package sender

import (
	"coldheater/internal/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)



func SendColdEmail(sender database.Sender, recipient string, content EmailContent) error {

	err := godotenv.Overload()
	if err != nil {
		return err
	}
	
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	//senderRef := sender.Name + " <" + sender.Email + ">"
	senderRef := sender.Email

	params := &resend.SendEmailRequest{
        From:    	 senderRef,
        To:      	 []string{recipient},
        Html:    	 content.Body,
        Subject: 	 content.Subject,
    }

    _, err = client.Emails.Send(params)

	return err
}