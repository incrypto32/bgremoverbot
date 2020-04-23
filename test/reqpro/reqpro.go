package reqpro

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func NewUploadRequest(url string, params map[string]string, paramName, path string, apiKey string) (*http.Request, error) {
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
	req, err := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("X-Api-Key", apiKey)
	req.Header.Add("User-Agent", "incrypto-telegram-bot by Krishnanand")
	return req, err
}
func Blah() {
	fmt.Println("blah...!!")
}
