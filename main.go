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

func main() {
	fmt.Println("APIKEY : ", apiKey)
	fmt.Println("BOT_TOKEN : ", botToken)
	fmt.Println("PORT : ", port)
	_ = http.ListenAndServe(":"+port, http.HandlerFunc(webHookHandler))
}

func webHookHandler(resp http.ResponseWriter, req *http.Request) {

	fmt.Println("________webhookHandler called________")
	// initializing a webhook request body type
	body := &webHookReqBody{}
	var funcName string = "webHookHandler : "
	fmt.Println(funcName, "webHook request body created")

	// error handling
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println(funcName, "could not decode request body", err)
		return
	}

<<<<<<< HEAD
	fmt.Println("Decoded request successfully :webHookHandler")
=======
	fmt.Println(funcName, "Decoded request successfully")
>>>>>>> tester

	if strings.Contains(strings.ToLower(body.Message.Text), "hello") {
		fmt.Println(funcName, "Message text contains hello")
		if err := sendReply(body.Message.Chat.ID, "Hi,There"); err != nil {
			log.Fatal(err)
			fmt.Println(funcName, "error in sending reply ", err)
			return
		}
	}

	fmt.Println(funcName, "Checked message text for hello")

	if len(body.Message.PhotoArr) != 0 {
		len := len(body.Message.PhotoArr)
		imageHandler(body.Message.PhotoArr[len-1].ID, body.Message.Chat.ID)

	}

	fmt.Println("reply sent")

}

// when the message is an image this function is called

func imageHandler(fileId string, chatId int64) {
	var funcName string = "imageHandler : "
	fmt.Println()
	fmt.Println()
	fmt.Println(funcName, "The Message is an image")

	imageUrl, err := getPhotoUrl(fileId)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(funcName, "getPhotoUrl Completed successfully")

	fmt.Println(funcName, "The image Url is :"+imageUrl)

	if imageUrl != "" {

		fmt.Println(funcName, "Image Url is not empty string ")

		imageUrl = botFileUrl + imageUrl

		fmt.Println(funcName, "created download url from image path ")

		// a := map[string]string{
		// 	"size":      "auto",
		// 	"image_url": imageUrl,
		// }

		fmt.Println(funcName, "map of params created ")

		// imgReq, err := rmvbgapi.NewUrlRequest(
		// 	"https://api.remove.bg/v1.0/removebg",
		// 	imageUrl,
		// 	a,
		// 	apiKey,
		// )
		fmt.Println(funcName, "Got img Req")
		if err != nil {
			fmt.Println("oh shit an error occured on imgUrl req")
			log.Fatal(funcName, err)
			return
		}

		fmt.Println(funcName, "img url request performed without errors")

		// fmt.Println(imgReq)
		resp, err := rmvbgapi.Driver("https://api.remove.bg/v1.0/removebg", apiKey, imageUrl)

		if err != nil {
			log.Fatal(funcName, err)
			return
		}
		if resp.StatusCode != http.StatusOK {
			fmt.Println(funcName, "Status Not OK")
			err = sendReply(chatId, "Oh sed lyf . Could Not Process Your request Now")
			if err != nil {
				log.Fatal(funcName, err)
			}
		} else {
			err = sendPhoto(chatId, resp)
			if err != nil {
				log.Fatal(funcName, err)
				return
			}
		}

	}
}

//This function is used to get the url of a photo from file_id
func getPhotoUrl(id string) (string, error) {
	var url string
	var funcName string = "getPhotoUrl : "
	fmt.Printf("\n\n")
	fmt.Println(funcName, "getPhotoUrl called  ")

	fileBody := &Result{}
	fmt.Println(funcName, "Tempelate for result created  ")
	resp, err := http.Get(
		botUrl + "getFile?file_id=" + id,
	)
	if err != nil {
		log.Fatal(funcName, err)
		return url, err
	}
	defer resp.Body.Close()
	fmt.Println(funcName, "http Get called from tgram api to get url of the photo  ")
	if err != nil {
		log.Fatal(funcName, "Shit an error occured while http Get")
		return url, err
	}
	fmt.Println(funcName, "no errors in http get  ")
	if resp.StatusCode != http.StatusOK {
		fmt.Println(funcName, errors.New("unexpected status"+resp.Status))
		return url, err
	}
	if err := json.NewDecoder(resp.Body).Decode(fileBody); err != nil {
		fmt.Println(funcName, "An error occured while decoding response into file Body ")
		return url, err
	}

	fmt.Println(funcName, "Successfully decoded response body into fileBody :getPhotoUrl")

	url = fileBody.File.Path
	fmt.Println(funcName, "The url is ", url)
	fmt.Printf("\n\n")
	return url, err
}

//Replies to a hello
func sendReply(chatID int64, msg string) error {
	var funcName string = "sendReply : "
	fmt.Println("\n\n", funcName, "Say polo called ")
	reqBody := &sendMessageReqBody{
		ChatID: chatID,
		Text:   msg,
	}
	fmt.Println(funcName, "request body tempelate created")
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Println(funcName, "An error occured while json Marshal")
		return err
	}
	fmt.Println(funcName, "json Marshal successfull :sayPolo")
	fmt.Println(funcName, "https://api.telegram.org/bot"+botToken+"/"+"sendMessage")
	resp, err := http.Post(
		"https://api.telegram.org/bot"+botToken+"/"+"sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		log.Fatal(funcName, err)
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + resp.Status)
	}
	fmt.Printf("\n\n")
	return nil
}

//Function to send a photo
func sendPhoto(chatID int64, inpresp *http.Response) error {
	var req *http.Request
	var err error
	var funcName string = "sendPhoto : "
	fmt.Println(funcName, "sendPhoto Called")
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
		log.Fatal(funcName, err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	fmt.Println(funcName, resp)
	fmt.Printf("\n\n")
	return err

}
