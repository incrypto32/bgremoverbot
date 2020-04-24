package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/incrypt0/bg_remover_bot/rmvbgapi"
)

//The webHook request template
type webHookReqBody struct {
	Message struct {
		Text string `json:"text"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		PhotoArr []Photo `json:"photo"`
	} `json:"message"`
}

//The Photo tempelate
type Photo struct {
	ID string `json:"file_id"`
}
type Result struct {
	File struct {
		Path string `json:"file_path"`
	} `json:"result"`
}

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

var apiKey string = os.Getenv("REMOVE_BG_API_KEY")
var botToken string = os.Getenv("TELEGRAM_TOKEN")
var port string = os.Getenv("PORT")
var botFileUrl = "https://api.telegram.org/file/bot" + botToken + "/"
var botUrl = "https://api.telegram.org/bot" + botToken + "/"

func webHookHandler(resp http.ResponseWriter, req *http.Request) {

	fmt.Println("________webhookHandler called________")
	// initializing a webhook request body type
	body := &webHookReqBody{}

	fmt.Println("webHook request body created :webHookHandler")

	// error handling
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body :webHookHandler", err)
		return
	}

	fmt.Println("Decoded request successfully")

	if strings.Contains(strings.ToLower(body.Message.Text), "hello") {
		fmt.Println("Message text contains hello :webHookHandler")
		if err := sendReply(body.Message.Chat.ID); err != nil {
			log.Fatal(err)
			fmt.Println("error in sending reply : :webHookHandler:", err)
			return
		}
	}

	fmt.Println("Checked message text for hello :webHookHandler")

	if len(body.Message.PhotoArr) != 0 {
		imageHandler(body.Message.PhotoArr[2].ID, body.Message.Chat.ID)

	}

	fmt.Println("reply sent")

}

// when the message is an image this function is called

func imageHandler(fileId string, chatId int64) {
	fmt.Println("The Message is an image :webHookHandler")

	imageUrl, err := getPhotoUrl(fileId)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("getPhotoUrl Completed successfully :webHookHandler")

	fmt.Println("The image Url is : ::webHookHandler:" + imageUrl)

	if imageUrl != "" {

		fmt.Println("checked whether image url is an emty string and its not :webHookHandler")

		imageUrl = botFileUrl + imageUrl

		fmt.Println("created download url from image path :webHookHandler :", imageUrl)

		a := map[string]string{
			"size":      "auto",
			"image_url": imageUrl,
		}

		fmt.Println("map of params created :webHookHandler :", a)

		imgReq, err := rmvbgapi.NewUrlRequest("https://api.remove.bg/v1.0/removebg", imageUrl, a, apiKey)
		fmt.Println("Got img Req")
		if err != nil {
			fmt.Println("oh shit an error occured on imgUrl req")
			log.Fatal(err)
			return
		}

		fmt.Println("img url request performed without errors :webHookHandler")

		fmt.Println(imgReq)
		resp, err := rmvbgapi.Driver("https://api.remove.bg/v1.0/removebg", apiKey, imageUrl)

		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(*resp)
		err = sendPhoto(chatId, resp)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

//This function is used to get the url of a photo from file_id
func getPhotoUrl(id string) (string, error) {
	var url string
	fmt.Println("getPhotoUrl called :getPhotoUrl :id is :", id)

	fileBody := &Result{}
	fmt.Println("Tempelate for result created :getPhotoUrl ")
	resp, err := http.Get(
		botUrl + "getFile?file_id=" + id,
	)
	if err != nil {
		log.Println(err)
		return url, err
	}
	defer resp.Body.Close()
	fmt.Println("http Get called from tgram api to get url of the photo :getPhotoUrl ")
	if err != nil {
		fmt.Println("Shit an error occured while http Get")
		return url, err
	}
	fmt.Println("no errors in http get :getPhotoUrl ")
	if resp.StatusCode != http.StatusOK {
		fmt.Println(errors.New("unexpected status ::getPhotoUrl:" + resp.Status))
		return url, err
	}
	if err := json.NewDecoder(resp.Body).Decode(fileBody); err != nil {
		fmt.Println("An error occured while decoding response into file Body :getPhotoUrl")
		return url, err
	}

	fmt.Println("Successfully decoded response body into fileBody :getPhotoUrl")

	url = fileBody.File.Path
	fmt.Println("The url is ::getPhotoUrl:", url)
	return url, err
}

//Replies to a hello
func sendReply(chatID int64) error {
	fmt.Println("Say polo called : sayPolo")
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   "Hi ,there",
	}
	fmt.Println("request body tempelate created : sayPolo")
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println("An error occured while json Marshal : sayPolo")
		return err
	}
	fmt.Println("json Marshal successfull :sayPolo")
	fmt.Println("https://api.telegram.org/bot" + botToken + "/" + "sendMessage")
	resp, err := http.Post(
		"https://api.telegram.org/bot"+botToken+"/"+"sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + resp.Status)
	}
	return nil
}

//Function to send a photo
func sendPhoto(chatID int64, inpresp *http.Response) error {
	var req *http.Request
	var err error
	fmt.Println("sendPhoto Called")
	fmt.Println()
	chatIDstr := strconv.FormatInt(chatID, 10)
	params := map[string]string{
		"chat_id": chatIDstr,
		"caption": "Haha",
	}
	fmt.Println(params)
	if inpresp != nil && inpresp.StatusCode != http.StatusOK {
		req, err = rmvbgapi.NewUploadRequest(botUrl+"sendDocument", params, "document", "Blah.png", "", nil)
	} else {
		fmt.Println(inpresp)
		req, err = rmvbgapi.NewUploadRequest(botUrl+"sendDocument", params, "document", "", "", inpresp)
	}
	if err != nil {
		log.Fatal(err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(resp)
	return err

}

func main() {
	fmt.Println(apiKey, botToken)
	_ = http.ListenAndServe(":"+port, http.HandlerFunc(webHookHandler))
}
