package pkg

import (
	"62-gpt/app/interfaces"
	"62-gpt/config"
	"62-gpt/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type WhatsappService struct {
	Config config.Config
}

func NewWhatsapp() interfaces.WhatsappPackage {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
	}

	return &WhatsappService{
		Config: loadedConfig,
	}
}

func (w WhatsappService) SendMessage(messageId string, from string, phoneNumberId string, answer string) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"context": map[string]string{
			"message_id": messageId,
		},
		"to": from,
		"text": map[string]string{
			"body": answer,
		},
		"type": "text",
	}

	utils.Fetch("POST", "https://graph.facebook.com/v17.0/"+phoneNumberId+"/messages", body, w.Config.WhatsappToken)

}

func (w WhatsappService) ReadMessage(messageId string, phoneNumberId string) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"status":            "read",
		"message_id":        messageId,
	}

	utils.Fetch("POST", "https://graph.facebook.com/v17.0/"+phoneNumberId+"/messages", body, w.Config.WhatsappToken)
}

func (w WhatsappService) UploadImage(phoneNumberId string, path string) interfaces.ResponseMedia {
	// Open the image file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a buffer to store the form data
	body := &bytes.Buffer{}
	writer := utils.NewWriter(body)

	// Create a form field for the image file
	imageField, err := writer.CreateFormFile("file", filepath.Base(file.Name()), "image/jpeg")
	if err != nil {
		log.Fatal(err)
	}

	// Copy the image data to the form field
	_, err = io.Copy(imageField, file)
	if err != nil {
		log.Fatal(err)
	}

	// Add the messaging_product field
	err = writer.WriteField("messaging_product", "whatsapp")
	if err != nil {
		log.Fatal(err)
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Create a POST request to the server
	request, err := http.NewRequest("POST", "https://graph.facebook.com/v17.0/"+phoneNumberId+"/media", body)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "image/png")
	request.Header.Set("Authorization", "Bearer "+w.Config.WhatsappToken)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Accept", "*/*")

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data interfaces.ResponseMedia
	err = json.Unmarshal(responseBody, &data)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return data
}

func (w WhatsappService) DownloadImage(phoneNumber string, id string) (string, error) {
	res, err := utils.Fetch("GET", "https://graph.facebook.com/v17.0/"+id+"?phone_number_id="+phoneNumber, nil, w.Config.WhatsappToken)
	if err != nil {
		return "", nil
	}

	attachment := interfaces.ImageResponse{
		Id:               res["id"].(string),
		FileSize:         res["file_size"].(float64),
		MessagingProduct: res["messaging_product"].(string),
		MimeType:         res["mime_type"].(string),
		Sha256:           res["sha256"].(string),
		Url:              res["url"].(string),
	}

	filename := utils.GenerateFileName(".jpg")
	err = utils.ImageDownload(attachment.Url, "tmp/img/"+filename, w.Config.WhatsappToken)
	if err != nil {
		return "", err
	}

	return "tmp/img/" + filename, nil
}

func (w WhatsappService) SendImageMessage(messageId string, from string, phoneNumberId string, id string, caption string) {
	body := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"context": map[string]string{
			"message_id": messageId,
		},
		"to":   from,
		"type": "image",
		"image": map[string]string{
			"id":      id,
			"caption": caption,
		},
	}

	utils.Fetch("POST", "https://graph.facebook.com/v17.0/"+phoneNumberId+"/messages", body, w.Config.WhatsappToken)
}
