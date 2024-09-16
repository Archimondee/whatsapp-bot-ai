package interfaces

type GptService interface {
	GPT3Response(string) (string, error)
	ConversationMessage(phone_number string, text string, message string) (response string, err error)
	EndConversation(phone_number string) error
	StartConversation(phone_number string, bot_number string) (response Message, err error)
}

type Message struct {
	PhoneNumber string `json:"phone_number"`
	Text        string `json:"text"`
}

type Whitelist struct {
	PhoneNumber string `json:"phone_number"`
}
