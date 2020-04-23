package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func NewUploadRequest(providerUrl string, params map[string]string, paramName, path string, apiKey string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	_, err = io.Copy(part, file)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", providerUrl, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("User-Agent", "incrypto-telegram-bot by Krishnanand")
	return req, err
}
func NewUrlRequest(providerUrl string, params map[string]string, apiKey string) (*http.Request, error) {
	fmt.Println("Called NewUrlRequest", params, providerUrl, apiKey)
	// We are setting up the key value pairs for the post method
	form := url.Values{}

	// here we add the key value pairs  to the Values tempelate
	for key, val := range params {
		form.Add(key, val)
	}

	// Here we are filling the request model
	req, err := http.NewRequest("POST", providerUrl, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// Adding headers and stuff
	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("User-Agent", "incrypto-telegram-bot by Krishnanand")
	return req, err
}
