package sender

import "coldheater/internal/database"

type EmailSender interface {
	SendColdEmail(sender database.Sender, recipient string, content EmailContent) error
	//SendWarmupReply(conversation Conversation, replyContent string) error
}

type ConversationManager interface {
	//GetPendingConversations(senderId string) ([]Conversation, error)
	//UpdateConversationState(conversaitonId string, newState ConversationState) error
}
