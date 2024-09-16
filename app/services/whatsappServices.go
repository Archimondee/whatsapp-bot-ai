package services

import (
	"62-gpt/app/interfaces"
	"62-gpt/config"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type WhatsappService struct {
	Config config.Config
}

func NewWhatsappService() interfaces.WhatsappAdapter {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)

	}
	return &WhatsappService{
		Config: loadedConfig,
	}
}

func (w *WhatsappService) GetWhatsappWebhook(ctx *gin.Context) (string, error) {
	verifyToken := w.Config.VerifyToken
	mode := ctx.Query("hub.mode")
	token := ctx.Query("hub.verify_token")
	challenge := ctx.Query("hub.challenge")

	if mode != "" && token != "" {
		if mode == "subscribe" && token == verifyToken {
			fmt.Println("WEBHOOK_VERIFIED")

			return challenge, nil
		} else {
			return "", errors.New("Forbidden")
		}
	} else {
		return "", errors.New("Forbidden")
	}
}

func (w *WhatsappService) PostWhatsappWebhook(ctx *gin.Context) (*interfaces.WhatsappResponse, *interfaces.WhatsappResponseImage, error) {
	var body map[string]interface{}
	err := ctx.BindJSON(&body)
	if err != nil {
		return nil, nil, err
	}

	//fmt.Printf("%+v\n", body)
	if body["object"] != nil {
		if entry, ok := body["entry"]; ok {
			if entries, ok := entry.([]interface{}); ok && len(entries) > 0 {
				if changes, ok := entries[0].(map[string]interface{})["changes"]; ok {
					if changesSlice, ok := changes.([]interface{}); ok && len(changesSlice) > 0 {
						if messages, ok := changesSlice[0].(map[string]interface{})["value"].(map[string]interface{})["messages"]; ok {
							if messagesSlice, ok := messages.([]interface{}); ok && len(messagesSlice) > 0 {
								messageID := messagesSlice[0].(map[string]interface{})["id"].(string)
								phoneNumberID := changesSlice[0].(map[string]interface{})["value"].(map[string]interface{})["metadata"].(map[string]interface{})["phone_number_id"].(string)
								from := messagesSlice[0].(map[string]interface{})["from"].(string)
								timestamp := messagesSlice[0].(map[string]interface{})["from"].(string)
								if msgBodyMap, ok := messagesSlice[0].(map[string]interface{})["text"].(map[string]interface{}); ok {
									if msgBody, ok := msgBodyMap["body"].(string); ok {
										if phoneNumberID != "" && from != "" && msgBody != "" && messageID != "" {
											response := &interfaces.WhatsappResponse{
												PhoneNumberId: phoneNumberID,
												From:          from,
												MessageBody:   msgBody,
												MessageId:     messageID,
												Timestamp:     timestamp,
											}
											return response, nil, nil
										}
										return nil, nil, errors.New("Listening")
									} else {
										return nil, nil, errors.New("Listening")
									}
								} else {
									if msgBodyMap, ok := messagesSlice[0].(map[string]interface{})["image"].(map[string]interface{}); ok {
										if caption, ok := msgBodyMap["caption"]; ok {
											image := interfaces.Image{
												ID:       msgBodyMap["id"].(string),
												Caption:  caption.(string),
												MimeType: msgBodyMap["mime_type"].(string),
												Sha256:   msgBodyMap["sha256"].(string),
											}
											response := &interfaces.WhatsappResponseImage{
												PhoneNumberId: phoneNumberID,
												From:          from,
												MessageId:     messageID,
												Image:         image,
												Timestamp:     timestamp,
											}
											return nil, response, nil
										} else {
											return nil, nil, errors.New("Listening")
										}
									}
								}
							}
						}
					}
				}
			}
		}
		return nil, nil, errors.New("Listening")
	}
	return nil, nil, errors.New("Listening")
}
