package pkg

import (
	"62-gpt/app/interfaces"
	"62-gpt/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var url = "https://cloud.leonardo.ai/api/rest/v1/generations"
var uploadUrl = "https://cloud.leonardo.ai/api/rest/v1/init-image"

type Client struct {
	Authorization     string
	Model             *interfaces.Model
	Request           *interfaces.Request
	Response          *interfaces.Response
	ResponseGenerate  *interfaces.ResponseGenerate
	ResponseInitImage *interfaces.ResponseInitImage
}

func NewClient(auth string, model *interfaces.Model) *Client {
	return &Client{
		Authorization: auth,
		Model:         model,
		Request:       model.CreateRequest(),
	}
}

func (c *Client) Create() (err error) {
	bodyInput, err := json.Marshal(c.Request.Input)
	if err != nil {
		return err
	}
	jsonString := string(bodyInput)
	payload := strings.NewReader(jsonString)
	//payload := strings.NewReader("{\"prompt\":\"An oil painting of a cat\",\"modelId\":\"6bef9f1b-29cb-40c7-b9df-32b51c1f67d3\",\"width\":1024,\"height\":1024,\"sd_version\":\"v2\",\"num_images\":1}")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		//log.Println(log.LogLevelError, err)
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.Authorization)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &c.Response)
	return
}

func (c *Client) Get(predictionId string) (err error) {
	// create a HTTP GET request
	req, err := http.NewRequest("GET", url+"/"+predictionId, nil)
	req.Header.Add("Authorization", "Bearer "+c.Authorization)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	//httpClient := http.Client{}
	// create a HTTP client and use it to send the request
	ok := make(chan bool)
	go func() {
		for {
			// send a HTTP GET request to the API
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				break
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)

			if err != nil {
				break
			}
			err = json.Unmarshal(body, &c.ResponseGenerate)

			if err != nil {
				break
			}

			if c.ResponseGenerate.GenerationsByPK.Status == "COMPLETE" {
				ok <- true
				break
			}
			if c.ResponseGenerate.GenerationsByPK.Status == "FAILED" {
				ok <- true
				break
			}
		}
	}()
	<-ok
	return
}

func (c *Client) UploadInitImage(imagePath string) (string, error) {
	payload := strings.NewReader("{\"extension\":\"jpeg\"}")
	req, _ := http.NewRequest("POST", uploadUrl, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+c.Authorization)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var response interfaces.ResponseInitImage
	err := json.Unmarshal(body, &response)
	if err != nil {
		//log.Println(log.LogLevelInfo, "err"+err.Error())
		return "", err
	}

	file, err := os.Open(imagePath)
	if err != nil {
		//log.Println(log.LogLevelInfo, "err1"+err.Error())
		return "", err
	}
	defer file.Close()

	bodyBytes := new(bytes.Buffer)
	writer := multipart.NewWriter(bodyBytes)
	var fieldData interfaces.FieldData
	errJson := json.Unmarshal([]byte(response.UploadInitImage.Fields), &fieldData)
	if errJson != nil {

		//log.Println(log.LogLevelInfo, "err3"+errJson.Error())
	}
	writer.WriteField("Content-Type", fieldData.ContentType)
	writer.WriteField("bucket", fieldData.Bucket)
	writer.WriteField("X-Amz-Algorithm", fieldData.XAmzAlgorithm)
	writer.WriteField("X-Amz-Credential", fieldData.XAmzCredential)
	writer.WriteField("X-Amz-Date", fieldData.XAmzDate)
	writer.WriteField("X-Amz-Security-Token", fieldData.XAmzSecurityToken)
	writer.WriteField("key", fieldData.Key)
	writer.WriteField("Policy", fieldData.Policy)
	writer.WriteField("X-Amz-Signature", fieldData.XAmzSignature)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		//log.Println(log.LogLevelInfo, "err2"+err.Error())
		return "", err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		//log.Println(log.LogLevelInfo, "err3"+err.Error())
		return "", err
	}

	err = writer.Close()
	if err != nil {
		//log.Println(log.LogLevelInfo, "err4"+err.Error())
		return "", err
	}

	request, err := http.NewRequest("POST", response.UploadInitImage.URL, bodyBytes)
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err = client.Do(request)
	if err != nil {
		return "", err
	}

	_, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return response.UploadInitImage.Id, err
}

func Dreamshaper6(prompt string, path string) (string, error) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}
	leo := interfaces.NewModel(prompt, loadedConfig.DreamshaperV6)
	leoClient := NewClient(loadedConfig.LeonardoApiKey, leo)
	if path != "" {
		id, err := leoClient.UploadInitImage(path)
		if err != nil {
			return "", err
		}
		if id == "" {
			return "", errors.New("Something when wrong")
		} else {
			leo.Input["init_image_id"] = id
			leo.Input["init_strength"] = 0.3
		}
	}
	err = leoClient.Create()
	if err != nil {
		//log.Println(log.LogLevelInfo, err.Error())
		return "", err
	}
	generationId := leoClient.Response.SdGenerationJob.GenerationId
	errGet := leoClient.Get(generationId)
	if errGet != nil {
		//log.Println(log.LogLevelInfo, "2"+err.Error())
		return "", errGet
	}

	return leoClient.ResponseGenerate.GenerationsByPK.GeneratedImages[0].URL, nil
}

func Creative(prompt string, path string) (string, error) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}
	leo := interfaces.NewModel(prompt, loadedConfig.LeonardoCreative)
	leoClient := NewClient(loadedConfig.LeonardoApiKey, leo)
	if path != "" {
		id, err := leoClient.UploadInitImage(path)
		if err != nil {
			return "", err
		}
		if id == "" {
			return "", errors.New("Something when wrong")
		} else {
			leo.Input["init_image_id"] = id
			leo.Input["init_strength"] = 0.3
		}
	}
	err = leoClient.Create()
	if err != nil {
		//log.Println(log.LogLevelInfo, err.Error())
		return "", err
	}
	generationId := leoClient.Response.SdGenerationJob.GenerationId
	errGet := leoClient.Get(generationId)
	if errGet != nil {
		//log.Println(log.LogLevelInfo, "2"+err.Error())
		return "", errGet
	}

	return leoClient.ResponseGenerate.GenerationsByPK.GeneratedImages[0].URL, nil
}

func Select(prompt string, path string) (string, error) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}
	leo := interfaces.NewModel(prompt, loadedConfig.LeonardoSelect)
	leoClient := NewClient(loadedConfig.LeonardoApiKey, leo)
	if path != "" {
		id, err := leoClient.UploadInitImage(path)
		if err != nil {
			return "", err
		}
		if id == "" {
			return "", errors.New("Something when wrong")
		} else {
			leo.Input["init_image_id"] = id
			leo.Input["init_strength"] = 0.3
		}
	}
	err = leoClient.Create()
	if err != nil {
		//log.Println(log.LogLevelInfo, err.Error())
		return "", err
	}
	generationId := leoClient.Response.SdGenerationJob.GenerationId
	errGet := leoClient.Get(generationId)
	if errGet != nil {
		//log.Println(log.LogLevelInfo, "2"+err.Error())
		return "", errGet
	}

	return leoClient.ResponseGenerate.GenerationsByPK.GeneratedImages[0].URL, nil
}

func Signature(prompt string, path string) (string, error) {
	loadedConfig, err := config.LoadConfig(".")
	if err != nil {
		return "", err
	}
	leo := interfaces.NewModel(prompt, loadedConfig.LeonardoSignature)
	leoClient := NewClient(loadedConfig.LeonardoApiKey, leo)
	if path != "" {
		id, err := leoClient.UploadInitImage(path)
		if err != nil {
			return "", err
		}
		if id == "" {
			return "", errors.New("Something when wrong")
		} else {
			leo.Input["init_image_id"] = id
			leo.Input["init_strength"] = 0.3
		}
	}
	err = leoClient.Create()
	if err != nil {
		//log.Println(log.LogLevelInfo, err.Error())
		return "", err
	}
	generationId := leoClient.Response.SdGenerationJob.GenerationId
	errGet := leoClient.Get(generationId)
	if errGet != nil {
		//log.Println(log.LogLevelInfo, "2"+err.Error())
		return "", errGet
	}

	return leoClient.ResponseGenerate.GenerationsByPK.GeneratedImages[0].URL, nil
}
