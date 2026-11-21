package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

var url string = "http://localhost:8080/"

type RequestConfig struct {
	URL     string
	Method  string
	Body    string
	Headers map[string]string
	Timeout time.Duration
}

func main() {
	// createChirp()
	//  createUser()
	reset()
}

func createChirp() {
	createUser := RequestConfig{
		URL:    url + "api/chirps",
		Method: http.MethodPost,
		Body: `
		{
			"body": "Hello, world!",
			"user_id": "bf7aafff-e289-44b7-b289-ef8aebee51f4"
		}
		`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Timeout: 10 * time.Second,
	}
	doRequest(createUser)
}

func createUser() {
	createUser := RequestConfig{
		URL:    url + "api/users",
		Method: http.MethodPost,
		Body:   "{\"email\": \"user@example.com\"}",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Timeout: 10 * time.Second,
	}
	doRequest(createUser)
}
func reset() {
	reset := RequestConfig{
		URL:    url + "admin/reset",
		Method: http.MethodPost,
		Body:   "",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Timeout: 10 * time.Second,
	}
	doRequest(reset)
}

func doRequest(c RequestConfig) {
	req, err := http.NewRequest(c.Method, c.URL, bytes.NewBufferString(c.Body))
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: c.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", resp.Status)
	fmt.Println("Body:", string(respBody))
}
