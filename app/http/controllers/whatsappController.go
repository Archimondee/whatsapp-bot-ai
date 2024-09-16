package controllers

import (
	"62-gpt/app/interfaces"
	"62-gpt/app/pkg"
	"62-gpt/config"
	"62-gpt/utils"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type WhatsappController struct {
	WhatsappService interfaces.WhatsappAdapter
	GptService      interfaces.GptService
	Whatsapp        interfaces.WhatsappPackage
}

func NewWhatsappController(whatsappService interfaces.WhatsappAdapter, gptService interfaces.GptService, whatsapp interfaces.WhatsappPackage) interfaces.WhatsappController {
	return &WhatsappController{
		WhatsappService: whatsappService,
		GptService:      gptService,
		Whatsapp:        whatsapp,
	}
}

func (ctrl *WhatsappController) GetWhatsappWebhook(ctx *gin.Context) {
	res, err := ctrl.WhatsappService.GetWhatsappWebhook(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ResponseData("error", err.Error(), nil))
		return
	}
	ctx.String(http.StatusOK, res)
	return
}

func (ctrl *WhatsappController) PostWhatsappWebhook(ctx *gin.Context) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		fmt.Println("Error", err)
	}
	text, image, err := ctrl.WhatsappService.PostWhatsappWebhook(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, utils.ResponseData("Listening", err.Error(), nil))
		return
	}

	if image != nil {
		err := CheckPhoneRecipientWhitelist(image.From)
		if err != nil {
			ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Your number not in whitelist")
			return
		}
		caption := image.Image.Caption
		ctrl.Whatsapp.ReadMessage(image.MessageId, image.PhoneNumberId)
		if caption != "" {
			path, err := ctrl.Whatsapp.DownloadImage(image.PhoneNumberId, image.Image.ID)
			if err != nil {
				ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
			}

			if strings.Contains(caption, loadedConfig.DreamshaperTag) {
				remaining := strings.TrimPrefix(caption, loadedConfig.DreamshaperTag)
				res, err := pkg.Dreamshaper6(remaining, path)
				if err != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				filename := utils.GenerateFileName(".jpg")
				resPath, errRes := utils.Download(res, "tmp/img", filename)
				if errRes != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				data := ctrl.Whatsapp.UploadImage(image.PhoneNumberId, resPath)
				ctrl.Whatsapp.SendImageMessage(image.MessageId, image.From, image.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)
				os.RemoveAll(path)
				os.RemoveAll(resPath)
			}

			if strings.Contains(caption, loadedConfig.CreativeTag) {
				remaining := strings.TrimPrefix(caption, loadedConfig.CreativeTag)
				res, err := pkg.Creative(remaining, path)
				if err != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				filename := utils.GenerateFileName(".jpg")
				resPath, errRes := utils.Download(res, "tmp/img", filename)
				if errRes != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				data := ctrl.Whatsapp.UploadImage(image.PhoneNumberId, resPath)
				ctrl.Whatsapp.SendImageMessage(image.MessageId, image.From, image.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)
				os.RemoveAll(path)
				os.RemoveAll(resPath)
			}

			if strings.Contains(caption, loadedConfig.SelectTag) {
				remaining := strings.TrimPrefix(caption, loadedConfig.SelectTag)
				res, err := pkg.Select(remaining, path)
				if err != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				filename := utils.GenerateFileName(".jpg")
				resPath, errRes := utils.Download(res, "tmp/img", filename)
				if errRes != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				data := ctrl.Whatsapp.UploadImage(image.PhoneNumberId, resPath)
				ctrl.Whatsapp.SendImageMessage(image.MessageId, image.From, image.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)
				os.RemoveAll(path)
				os.RemoveAll(resPath)
			}

			if strings.Contains(caption, loadedConfig.SignatureTag) {
				remaining := strings.TrimPrefix(caption, loadedConfig.SignatureTag)
				res, err := pkg.Signature(remaining, path)
				if err != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				filename := utils.GenerateFileName(".jpg")
				resPath, errRes := utils.Download(res, "tmp/img", filename)
				if errRes != nil {
					ctrl.Whatsapp.SendMessage(image.MessageId, image.From, image.PhoneNumberId, "Something is error in AI")
					return
				}

				data := ctrl.Whatsapp.UploadImage(image.PhoneNumberId, resPath)
				ctrl.Whatsapp.SendImageMessage(image.MessageId, image.From, image.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)
				os.RemoveAll(path)
				os.RemoveAll(resPath)
				return
			}
		}
	}

	if text != nil {
		ok, errAdd := AddAndRemovePhoneNumber(text.MessageBody, text.From)
		if errAdd != nil {
			return
		}
		if ok == "add" {
			ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Success add number.")
			return
		}
		if ok == "remove" {
			ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Success remove number.")
			return
		}

		err := CheckPhoneRecipientWhitelist(text.From)
		if err != nil {
			ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Your number not in whitelist")
			return
		}
		oneMinuteDuration := time.Minute
		timeA := time.Now()
		timeB, err := strconv.ParseInt(text.Timestamp, 10, 64)
		if err != nil {
			fmt.Println("Error parsing timestamp B:", err)
			return
		}

		timestampB := time.Unix(timeB, 0)
		timeDifference := timeA.Sub(timestampB)
		if timeDifference < oneMinuteDuration {
			ctrl.Whatsapp.ReadMessage(text.MessageId, text.PhoneNumberId)

			if strings.HasPrefix(text.MessageBody, loadedConfig.DreamshaperTag) ||
				strings.HasPrefix(text.MessageBody, loadedConfig.CreativeTag) ||
				strings.HasPrefix(text.MessageBody, loadedConfig.SignatureTag) ||
				strings.HasPrefix(text.MessageBody, loadedConfig.SelectTag) {
				if strings.Contains(text.MessageBody, loadedConfig.SignatureTag) {
					remaining := strings.TrimPrefix(text.MessageBody, loadedConfig.SignatureTag)
					res, err := pkg.Signature(remaining, "")
					if err != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					filename := utils.GenerateFileName(".jpg")
					resPath, errRes := utils.Download(res, "tmp/img", filename)
					if errRes != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					data := ctrl.Whatsapp.UploadImage(text.PhoneNumberId, resPath)
					ctrl.Whatsapp.SendImageMessage(text.MessageId, text.From, text.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)

					os.RemoveAll(resPath)
					return
				}
				if strings.Contains(text.MessageBody, loadedConfig.CreativeTag) {
					remaining := strings.TrimPrefix(text.MessageBody, loadedConfig.CreativeTag)
					res, err := pkg.Creative(remaining, "")
					if err != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					filename := utils.GenerateFileName(".jpg")
					resPath, errRes := utils.Download(res, "tmp/img", filename)
					if errRes != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					data := ctrl.Whatsapp.UploadImage(text.PhoneNumberId, resPath)
					ctrl.Whatsapp.SendImageMessage(text.MessageId, text.From, text.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)

					os.RemoveAll(resPath)
					return
				}
				if strings.Contains(text.MessageBody, loadedConfig.DreamshaperTag) {
					remaining := strings.TrimPrefix(text.MessageBody, loadedConfig.DreamshaperTag)
					res, err := pkg.Dreamshaper6(remaining, "")
					if err != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					filename := utils.GenerateFileName(".jpg")
					resPath, errRes := utils.Download(res, "tmp/img", filename)
					if errRes != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					data := ctrl.Whatsapp.UploadImage(text.PhoneNumberId, resPath)
					ctrl.Whatsapp.SendImageMessage(text.MessageId, text.From, text.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)

					os.RemoveAll(resPath)
					return
				}
				if strings.Contains(text.MessageBody, loadedConfig.SelectTag) {
					remaining := strings.TrimPrefix(text.MessageBody, loadedConfig.SelectTag)
					res, err := pkg.Select(remaining, "")
					if err != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					filename := utils.GenerateFileName(".jpg")
					resPath, errRes := utils.Download(res, "tmp/img", filename)
					if errRes != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Something is error in AI")
						return
					}

					data := ctrl.Whatsapp.UploadImage(text.PhoneNumberId, resPath)
					ctrl.Whatsapp.SendImageMessage(text.MessageId, text.From, text.PhoneNumberId, data.ID, "Image complete \nDownload : "+res)

					os.RemoveAll(resPath)
					return
				}
			} else {
				if strings.HasPrefix(text.MessageBody, loadedConfig.WhatsappDeleteConversation) ||
					strings.HasPrefix(text.MessageBody, loadedConfig.WhatsappNewConversation) {
					err := ctrl.GptService.EndConversation(text.From)
					if err == nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "New Conversation")
						return
					}
				} else {
					var data interfaces.Message

					result := utils.DB.Table("messages").First(&data).Where("phone_number = ?", text.From)
					if result.Error != nil {
						res, err := ctrl.GptService.StartConversation(text.From, text.PhoneNumberId)
						if err != nil {
							ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "Sorry AI can answer this time.")
							return
						}
						data = res
					}

					response, err := ctrl.GptService.ConversationMessage(text.From, text.MessageBody, data.Text)
					if err != nil {
						ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, err.Error())
						//ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, "GPT Token is full, you can try to /new or /add to create new conversation")
						return
					}
					ctrl.Whatsapp.SendMessage(text.MessageId, text.From, text.PhoneNumberId, response)
					return
				}
			}

		}
	}

	return
}

func AddAndRemovePhoneNumber(message string, from string) (string, error) {
	for _, value := range utils.SuperNumber {
		if value == from {
			if strings.HasPrefix(message, "/add") {
				re := regexp.MustCompile(`(\d+)`)
				match := re.FindString(message)
				if match != "" {
					var data = &interfaces.Whitelist{
						PhoneNumber: match,
					}
					response := utils.DB.Table("whitelist").Create(data)
					if response.Error != nil {
						return "", response.Error
					}
					if response.RowsAffected > 0 {
						return "add", nil
					} else {
						return "", errors.New("Error")
					}
				} else {
					return "", errors.New("Error")
				}
			}

			if strings.HasPrefix(message, "/remove") {
				re := regexp.MustCompile(`(\d+)`)
				match := re.FindString(message)
				if match != "" {
					response := utils.DB.Table("whitelist").Where("phone_number = ?", match).Delete(&interfaces.Whitelist{})
					if response.Error != nil {
						return "", response.Error
					}

					if response.RowsAffected > 0 {
						return "remove", nil
					} else {
						return "", errors.New("Error")
					}
				} else {
					return "", errors.New("Error")
				}

			}
			return "", nil
		}
	}
	return "", nil
}

func CheckPhoneRecipientWhitelist(from string) error {
	var whitelist interfaces.Message
	resp := utils.DB.Table("whitelist").First(&whitelist).Where("phone_number = ?", from)
	if resp.Error != nil {
		return resp.Error
	}

	return nil
}
