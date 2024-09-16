package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func Download(url string, folderPath string, customName string) (string, error) {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, 0755)
	}
	response, err := http.Get(url)
	//if bearer != "" {
	//response.Header.Add("Authorization", "Bearer "+bearer)
	//}

	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var fileName string
	if customName != "" {
		fileName = customName
	} else {
		fileName = path.Base(url)
	}
	filePath := filepath.Join(folderPath, fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func ImageDownload(url, filePath, token string) error {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	request.Header.Set("Authorization", "Bearer "+token)

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to download image, status code: %d", response.StatusCode)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
