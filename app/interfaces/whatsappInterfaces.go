package interfaces

import (
	"github.com/gin-gonic/gin"
)

type WhatsappAdapter interface {
	GetWhatsappWebhook(ctx *gin.Context) (string, error)
	PostWhatsappWebhook(ctx *gin.Context) (*WhatsappResponse, *WhatsappResponseImage, error)
}

type WhatsappController interface {
	GetWhatsappWebhook(ctx *gin.Context)
	PostWhatsappWebhook(ctx *gin.Context)
}

type WhatsappResponse struct {
	PhoneNumberId string `json:"phone_number_id"`
	From          string `json:"from"`
	MessageBody   string `json:"message_body"`
	MessageId     string `json:"message_id"`
	Timestamp     string `json:"timestamp"`
}

type WhatsappResponseImage struct {
	PhoneNumberId string `json:"phone_number_id"`
	From          string `json:"from"`
	MessageId     string `json:"message_id"`
	Image         Image  `json:"image"`
	Timestamp     string `json:"timestamp"`
}

type Image struct {
	Caption  string `json:"caption"`
	ID       string `json:"id"`
	MimeType string `json:"mime_type"`
	Sha256   string `json:"sha256"`
}

type ResponseMedia struct {
	ID string `json:"id"`
}

type WhatsappPackage interface {
	SendMessage(messageId string, from string, phoneNumberId string, answer string)
	ReadMessage(messageId string, phoneNumberId string)
	SendImageMessage(messageId string, from string, phoneNumberId string, id string, caption string)
	UploadImage(phoneNumberId string, path string) ResponseMedia
	DownloadImage(phoneNumber string, id string) (string, error)
}

type ImageResponse struct {
	FileSize         float64 `json:"file_size"`
	Id               string  `json:"id"`
	MessagingProduct string  `json:"messaging_product"`
	MimeType         string  `json:"mime_type"`
	Sha256           string  `json:"sha256"`
	Url              string  `json:"url"`
}
