package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId      int                 `json:"chat_id"`
	Text        string              `json:"text"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

type KeyboardButton struct {
	Text map[string]string `json:"text"`
}

var StartKeyboard = [][]KeyboardButton{
	{{map[string]string{"text": "Видео"}}},
	{{map[string]string{"text": "Музыка"}}},
	{{map[string]string{"text": "Мемы"}}},
	{{map[string]string{"text": "Соцсети"}}},
}

func main() {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Println("TELEGRAM_BOT_TOKEN environment variable does not exist")
		return
	}

	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Something went wrong: ", err.Error())
		}
		for _, update := range updates {
			err = respond(botUrl, update)
			offset = update.UpdateId + 1
		}

		fmt.Println(updates)
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

func respond(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId

	switch update.Message.Text {
	case "/start":
		handleStart(&botMessage, &update)
	default:
		handleDefault(&botMessage, &update)
	}

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}

	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}

func handleStart(botMessage *BotMessage, update *Update) {
	botMessage.Text = "Привет, это бот для фанатов Александра Пушного, тут можно найти разную инфу, видео, мемы, песни."
	botMessage.ReplyMarkup.Keyboard = StartKeyboard
}

func handleDefault(botMessage *BotMessage, update *Update) {
	log.Println("Unknown message: " + update.Message.Text)
	botMessage.Text = "Извини, я не знаю такой команды"
}
