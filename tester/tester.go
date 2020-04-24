package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/incrypt0/bg_remover_bot/rmvbgapi"
)

func main() {
	const imageUrl = "https://i2.wp.com/digital-photography-school.com/wp-content/uploads/2016/06/300.jpg?resize=750%2C499&ssl=1"
	var apiKey = os.Getenv("REMOVE_BG_API_KEY")
	// var apiKey = ""
	f, _ := os.Create("a.png")
	a := map[string]string{
		"size":      "auto",
		"image_url": imageUrl,
	}
	req, err := rmvbgapi.NewUrlRequest("https://api.remove.bg/v1.0/removebg", imageUrl, a, apiKey)
	if err != nil {
		log.Fatal(err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	resp.Body.Close()
	f.Close()
}
