package rmvbgapi

import (
	"log"
	"net/http"
	"os"
)

// func main() {
// 	//Here's our apikey
// 	const apiKey = "76b1Xqw9i82pTVtBC6YFQnSk"
// 	Driver(apiKey)
// 	// fmt.Println(path)

// }
func Driver(providerUrl string, apiKey string, imageUrl string) (*http.Response, error) {
	req, err := urlRequestGen(providerUrl, apiKey, imageUrl)
	if err != nil {
		log.Fatal(err)
	}

	//creating an http client
	client := &http.Client{}

	//client performing  the generated request template
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	} else {

		// //creating a byte buffer to read from the http response recieved
		// body := &bytes.Buffer{}

		// out, err := os.Create("Blah.png")
		// if err != nil {
		// 	// panic?
		// }
		// defer out.Close()
		// io.Copy(out, resp.Body)
		// //using the buffer and setting it up to read from body of http Response
		// _, err = body.ReadFrom(resp.Body)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// //Closing the response body
		// resp.Body.Close()

	}
	return resp, err
}

// upload file request geneartor
func UploadFileRequestGen(providerUrl string, apiKey string, pathToFile string) (*http.Request, error) {
	//Getting the path of the image file to be uploaded

	path, _ := os.Getwd()
	path += "/a.png"
	//extraParameters for the http request
	extraParams := map[string]string{
		"size":       "auto",
		"image_file": "blah",
	}

	// Function call which generates the request tempelate
	req, err := NewUploadRequest(providerUrl, extraParams, "image_file", path, apiKey, nil)
	return req, err
}

// url request genearator
func urlRequestGen(providerUrl string, apiKey string, imageUrl string) (*http.Request, error) {
	extraParams := map[string]string{
		"size":      "auto",
		"image_url": imageUrl,
	}
	req, err := NewUrlRequest(
		providerUrl,
		imageUrl,
		extraParams,
		apiKey,
	)
	return req, err
}
