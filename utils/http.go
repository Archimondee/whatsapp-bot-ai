package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func Fetch(method string, url string, body map[string]interface{}, token string) (res map[string]interface{}, err error) {
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err.Error())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err := ParseHTTPError(resp.Body)
		return res, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	err = json.Unmarshal(bodyBytes, &req)

	var b bytes.Buffer
	_, err = io.Copy(&b, resp.Body)

	if err != nil {
		return res, err
	}
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&res)

	if err != nil {
		return res, err
	}

	return res, err
}

func ParseHTTPError(body io.Reader) (err error) {
	var errRes map[string]map[string]interface{}
	err = json.NewDecoder(body).Decode(&errRes)
	if err != nil {
		return fmt.Errorf("unparsed error message")
	}
	msg := fmt.Sprintf("%s", errRes["error"]["message"])
	return errors.New(msg)
}
