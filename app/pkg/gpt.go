package pkg

import (
	"62-gpt/app/interfaces"
	"62-gpt/config"
	"62-gpt/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	OpenAI "github.com/sashabaranov/go-openai"
)

var OAIClient *OpenAI.Client

type GptService struct {
	Config config.Config
}

func NewGptService() interfaces.GptService {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Printf("cannot load loadedConfig: %w", err)
	}
	OAIClient = OpenAI.NewClient(loadedConfig.OpenAIApiKey)
	return &GptService{
		Config: loadedConfig,
	}
}

func (g *GptService) GPT3Response(question string) (string, error) {
	var gptResponseText string

	token, _ := strconv.ParseInt(g.Config.OpenAIGptModelToken, 0, 0)
	temperature, _ := strconv.ParseFloat(g.Config.OpenAIGptModelTemperature, 32)
	modelTopP, _ := strconv.ParseFloat(g.Config.OpenAIGptModelTopP, 32)
	penaltyPresence, _ := strconv.ParseFloat(g.Config.OpenAIGptModelPenaltyFrequency, 32)
	penaltyFrequency, _ := strconv.ParseFloat(g.Config.OpenAIGptModelPenaltyFrequency, 32)

	gptChatMode := regexp.MustCompile("\\b(?i)(" + "gpt-3\\.5" + ")")
	if gptChatMode.MatchString(g.Config.OpenAIGptModelName) {
		gptRequest := OpenAI.ChatCompletionRequest{
			Model:            g.Config.OpenAIGptModelName,
			MaxTokens:        int(token),
			Temperature:      float32(temperature),
			TopP:             float32(modelTopP),
			PresencePenalty:  float32(penaltyPresence),
			FrequencyPenalty: float32(penaltyFrequency),
			Messages: []OpenAI.ChatCompletionMessage{
				{
					Role:    OpenAI.ChatMessageRoleUser,
					Content: question,
				},
			},
		}

		gptResponse, err := OAIClient.CreateChatCompletion(
			context.Background(),
			gptRequest,
		)

		if err != nil {
			return "", err
		}

		if len(gptResponse.Choices) > 0 {
			gptResponseText = gptResponse.Choices[0].Message.Content
		}
	} else {
		gptRequest := OpenAI.CompletionRequest{
			Model:            g.Config.OpenAIGptModelName,
			MaxTokens:        int(token),
			Temperature:      float32(temperature),
			TopP:             float32(modelTopP),
			PresencePenalty:  float32(penaltyPresence),
			FrequencyPenalty: float32(penaltyFrequency),
			Prompt:           question,
		}

		gptResponse, err := OAIClient.CreateCompletion(
			context.Background(),
			gptRequest,
		)

		if err != nil {
			return "", err
		}

		if len(gptResponse.Choices) > 0 {
			gptResponseText = gptResponse.Choices[0].Text
		}
	}

	gptResponseBuffer := strings.TrimSpace(gptResponseText)
	gptResponseBuffer = strings.TrimLeft(gptResponseBuffer, "?\n")
	gptResponseBuffer = strings.TrimLeft(gptResponseBuffer, "!\n")
	gptResponseBuffer = strings.TrimLeft(gptResponseBuffer, ":\n")
	gptResponseBuffer = strings.TrimLeft(gptResponseBuffer, "'\n")
	gptResponseBuffer = strings.TrimLeft(gptResponseBuffer, ".\n")
	gptResponseBuffer = strings.TrimLeft(gptResponseBuffer, "\n")

	return gptResponseBuffer, nil
}

func (g *GptService) StartConversation(phone_number string, bot_number string) (response interfaces.Message, err error) {
	var data interfaces.Message
	token, _ := strconv.ParseInt(g.Config.OpenAIGptModelToken, 0, 0)
	temperature, _ := strconv.ParseFloat(g.Config.OpenAIGptModelTemperature, 32)
	modelTopP, _ := strconv.ParseFloat(g.Config.OpenAIGptModelTopP, 32)
	penaltyPresence, _ := strconv.ParseFloat(g.Config.OpenAIGptModelPenaltyFrequency, 32)
	penaltyFrequency, _ := strconv.ParseFloat(g.Config.OpenAIGptModelPenaltyFrequency, 32)
	var msg string
	for _, message := range utils.Chat {
		if message.PhoneNumberId == phone_number {
			msg = message.InitChat
		} else {
			msg = "You are a helpful user"
		}
	}

	gptChatMode := regexp.MustCompile("\\b(?i)(" + "gpt-3\\.5" + ")")
	if gptChatMode.MatchString(g.Config.OpenAIGptModelName) {

		req := OpenAI.ChatCompletionRequest{
			Model:            g.Config.OpenAIGptModelName,
			MaxTokens:        int(token),
			Temperature:      float32(temperature),
			TopP:             float32(modelTopP),
			PresencePenalty:  float32(penaltyPresence),
			FrequencyPenalty: float32(penaltyFrequency),
			Stream:           false,
			Messages: []OpenAI.ChatCompletionMessage{
				{
					Role:    "system",
					Content: msg,
				},
			},
		}
		jsonString, err := json.Marshal(req.Messages)
		if err != nil {
			return data, err
		}
		data = interfaces.Message{
			PhoneNumber: phone_number,
			Text:        string(jsonString),
		}
		result := utils.DB.Table("messages").Create(data)
		if result.Error != nil {
			return data, err
		}
	}

	return data, nil
}

func (g *GptService) EndConversation(phone_number string) error {
	err := utils.DB.Table("messages").Where("phone_number = ?", phone_number).Delete(&interfaces.Message{})
	if err != nil {
		return err.Error
	}

	return nil
}

func (g *GptService) ConversationMessage(phone_number string, text string, message string) (response string, err error) {
	var dataMessage []OpenAI.ChatCompletionMessage

	token, _ := strconv.ParseInt(g.Config.OpenAIGptModelToken, 0, 0)
	temperature, _ := strconv.ParseFloat(g.Config.OpenAIGptModelTemperature, 32)
	modelTopP, _ := strconv.ParseFloat(g.Config.OpenAIGptModelTopP, 32)
	penaltyPresence, _ := strconv.ParseFloat(g.Config.OpenAIGptModelPenaltyFrequency, 32)
	penaltyFrequency, _ := strconv.ParseFloat(g.Config.OpenAIGptModelPenaltyFrequency, 32)

	err = json.Unmarshal([]byte(message), &dataMessage)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	dataMessage = append(dataMessage, OpenAI.ChatCompletionMessage{
		Role:    OpenAI.ChatMessageRoleUser,
		Content: text,
	})

	req := OpenAI.ChatCompletionRequest{
		Model:            g.Config.OpenAIGptModelName,
		MaxTokens:        int(token),
		Temperature:      float32(temperature),
		TopP:             float32(modelTopP),
		PresencePenalty:  float32(penaltyPresence),
		FrequencyPenalty: float32(penaltyFrequency),
		Messages:         dataMessage,
	}
	//log.Println(log.LogLevelInfo, req)

	resp, err := OAIClient.CreateChatCompletion(context.Background(), req)
	fmt.Println("data", resp, err)
	if err != nil {
		return "", err
	}

	dataMessage = append(dataMessage, resp.Choices[0].Message)

	jsonString, err := json.Marshal(dataMessage)
	if err != nil {
		//log.Println(log.LogLevelInfo, err.Error())
	}
	var finalData = interfaces.Message{
		PhoneNumber: phone_number,
		Text:        string(jsonString),
	}

	result := utils.DB.Table("messages").Where("phone_number = ?", phone_number).Updates(&finalData)
	if result.Error != nil {
		return "", errors.New("Errors")
	} else {
		return resp.Choices[0].Message.Content, nil
	}
}
