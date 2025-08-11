package sender

type EmailSender interface {
	SendColdEmail(sender Sender, recipient string, content EmailContent) error
	SendWarmupReply(conversation Conversation, replyContent string) error
}

type ConversationManager interface {
	GetPendingConversations(senderId string) ([]Conversation, error)
	UpdateConversationState(conversaitonId string, newState ConversationState) error
}
