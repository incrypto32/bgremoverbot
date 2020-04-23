package rmvbgapi

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

func NewUploadRequest(providerUrl string, params map[string]string, paramName, path string, apiKey string, inresp *http.Response) (*http.Request, error) {
	var file *os.File

	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer file.Close()
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if path != "" {
		part, err := writer.CreateFormFile(paramName, filepath.Base(path))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		_, err = io.Copy(part, file)

	} else {
		part, err := writer.CreateFormFile(paramName, "myimage")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		_, err = io.Copy(part, inresp.Body)
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err := writer.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req, err := http.NewRequest("POST", providerUrl, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if apiKey != "" {
		req.Header.Add("X-Api-Key", apiKey)
	}
	req.Header.Add("User-Agent", "incrypto-telegram-bot by Krishnanand")
	return req, err
}
func NewUrlRequest(providerUrl string, imageUrl string, params map[string]string, apiKey string) (*http.Request, error) {

	fmt.Println("Called NewUrlRequest with provider url and api KEy:NewUrlRequest :", providerUrl, apiKey)
	fmt.Println("The image url is :NewUrlRequest: ", imageUrl)
	fmt.Println("The params are :NewUrlRequest", params)

	// We are setting up the key value pairs for the post method
	form := url.Values{}

	fmt.Println("Created the for tempelate for the http Request ::NewUrlRequest ")

	// here we add the key value pairs  to the Values tempelate
	for key, val := range params {
		form.Add(key, val)
		fmt.Println("Added :NewUrlRequest :", key, val)
	}
	fmt.Println("Added key value pairs : NewUrlRequest :")
	fmt.Println("Form is :NewUrlRequest:")
	// Here we are filling the request model

	req, err := http.NewRequest("POST", providerUrl, strings.NewReader(form.Encode()))

	if err != nil {
		fmt.Println("shit an error while creating new http req :NewUrlRequest")
		return nil, err
	}
	fmt.Println("Request created with no error  :NewUrlRequest  :")
	req.PostForm = form
	fmt.Println("Added form to req ::NewUrlRequest :")
	// Adding headers and stuff
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if apiKey != "" {
		fmt.Println("Checked if apikey is empty and yes it is :NewUrlRequest")
		req.Header.Add("X-Api-Key", apiKey)
	}
	fmt.Println("Everythings fine in here ::NewUrlRequest")
	fmt.Println(*req, ":NewUrlRequest")
	// req.Header.Add("User-Agent", "incrypto-telegram-bot by Krishnanand")
	return req, err
}
