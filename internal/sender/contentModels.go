package sender

type EmailContent struct {
	Subject string
	Body string
	Attachments []MailAttachment
}

type MailAttachment struct{
	Path string
	Filename string
}