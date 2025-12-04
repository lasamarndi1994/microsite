package service

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func SendWhatAppMessage() {
	url := "https://gate.whapi.cloud/messages/text?token=s7u9lV7qTvM7JHsSgS7UOblZpYa7NHMm"
	method := "POST"

	payload := strings.NewReader(`{
		"to": "9337243362",
		"body": "Hello, this message was sent via API!"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer s7u9lV7qTvM7JHsSgS7UOblZpYa7NHMm")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
