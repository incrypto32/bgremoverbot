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
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

	} else {
		part, err := writer.CreateFormFile(paramName, "myimage.png")
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		_, err = io.Copy(part, inresp.Body)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		defer inresp.Body.Close()
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
	req.Header.Add("User-Agent", "incrypto-telegram-bot by incrypt0")
	return req, err
}

func NewUrlRequest(
	providerUrl string,
	imageUrl string,
	params map[string]string,
	apiKey string,
) (*http.Request, error) {
	var funcName string = "NewUrlRequest : "
	fmt.Printf("\n\n")
	fmt.Println(funcName, "Called NewUrlRequest with provider url and api key")
	// fmt.Println(funcName, "The image url is  ", imageUrl)
	// fmt.Println(funcName, "The params are ", params)
	// We are setting up the key value pairs for the post method
	form := url.Values{}

	fmt.Println(funcName, "Created the for tempelate for the http Request ")

	// here we add the key value pairs  to the Values tempelate
	for key, val := range params {
		form.Add(key, val)
		fmt.Printf("\n")
		fmt.Println(funcName, "Added :", key, val)
	}
	fmt.Println()
	fmt.Println(funcName, "Added key value pairs  :")
	// fmt.Println("Form is :NewUrlRequest:")
	// Here we are filling the request model

	req, err := http.NewRequest("POST", providerUrl, strings.NewReader(form.Encode()))

	if err != nil {
		fmt.Println(funcName, "shit an error while creating new http req ")
		return nil, err
	}
	fmt.Println(funcName, "Request created with no error ")
	req.PostForm = form
	fmt.Println(funcName, "Added form to req:")
	// Adding headers and stuff
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if apiKey != "" {
		fmt.Println(funcName, "Checked if apikey is empty and yes it is ")
		req.Header.Add("X-Api-Key", apiKey)
	}
	fmt.Println(funcName, "Everythings fine in here ")
	// fmt.Println(funcName, *req)
	// req.Header.Add("User-Agent", "incrypto-telegram-bot by Krishnanand")
	fmt.Printf("\n\n")
	return req, err
}
