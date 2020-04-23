package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	//Here's our apikey
	const apiKey = "76b1Xqw9i82pTVtBC6YFQnSk"

	//Getting the path of the image file to be uploaded
	path, _ := os.Getwd()
	path += "/a.png"
	//extraParameters for the http request
	// extraParams := map[string]string{
	// 	"size":       "auto",
	// 	"image_file": "blah",
	// }

	//Function call which generates the request tempelate
	// req, err := reqpro.NewUploadRequest("https://api.remove.bg/v1.0/removebg", extraParams, "image_file", path, apiKey)

	extraParams := map[string]string{
		"size":      "auto",
		"image_url": "https://economictimes.indiatimes.com/thumb/msid-75059227,width-1200,height-900,resizemode-4,imgsize-247188/samir-racch.jpg?from=mdr",
	}
	req, err := NewUrlRequest(
		"https://api.remove.bg/v1.0/removebg",
		extraParams,
		apiKey,
	)
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

		//creating a byte buffer to read from the http response recieved
		body := &bytes.Buffer{}

		out, err := os.Create("Blah.png")
		if err != nil {
			// panic?
		}
		defer out.Close()
		io.Copy(out, resp.Body)
		//using the buffer and setting it up to read from body of http Response
		_, err = body.ReadFrom(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		//Closing the response body
		resp.Body.Close()
		fmt.Println(resp.StatusCode)
		fmt.Println(resp.Header)
		fmt.Println(body.Bytes())
	}
	// fmt.Println(path)

}
