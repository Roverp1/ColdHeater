package sender

import (
	"coldheater/internal/database"
	"os"

	"github.com/joho/godotenv"
	"github.com/resend/resend-go/v2"
)



func SendColdEmail(sender database.Sender, recipient string, content EmailContent) error {

	err := godotenv.Load()
	if err != nil {
		return err
	}
	
	apiKey := os.Getenv("RESEND_API_KEY")

	client := resend.NewClient(apiKey)

	//senderRef := sender.Name + " <" + sender.Email + ">"
	senderRef := sender.Email

	var atachments []*resend.Attachment;

	for _, mailAttachment := range content.Attachments{
		resendAtachment := &resend.Attachment{ 
			Path:     mailAttachment.Path,
    		Filename: mailAttachment.Filename,
		}
		atachments = append(atachments, resendAtachment)
	}	

	params := &resend.SendEmailRequest{
        From:    	 senderRef,
        To:      	 []string{recipient},
        Html:    	 content.Body,
        Subject: 	 content.Subject,
		Attachments: atachments,
    }

    _, err = client.Emails.Send(params)

	return err
}