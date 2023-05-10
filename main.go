package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/orijtech/gfycat/v1"
)

func main() {
	bot, err := NewGifBot()
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}

func initKeyboards() []tgbotapi.ReplyKeyboardMarkup {
	keyboards := make([]tgbotapi.ReplyKeyboardMarkup, 5)

	keyboards[0] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cats"),
			tgbotapi.NewKeyboardButton("Dogs"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Capybaras"),
			tgbotapi.NewKeyboardButton("Your Request"),
		),
	)
	keyboards[1] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Just Cats"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Cat memes"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)
	keyboards[2] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Just Dogs"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Dog memes"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)
	keyboards[3] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Just Capybaras"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Capybaras memes"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)
	keyboards[4] = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("I'm lucky"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("<- Back"),
		),
	)
	return keyboards
}

type GifBot struct {
	bot          *tgbotapi.BotAPI
	updateConfig tgbotapi.UpdateConfig
	keyboards    []tgbotapi.ReplyKeyboardMarkup
	gfycat       *gfycat.Client
}

func NewGifBot() (*GifBot, error) {
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	if botToken == "" {
		log.Println("TELEGRAM_BOT_TOKEN environment variable does not exist")
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable does not exist")
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	offset, err := strconv.ParseInt(os.Getenv("UPDATE_OFFSET"), 10, 32)
	if err != nil {
		log.Println("UPDATE_OFFSET environment variable does not exist")
		return nil, fmt.Errorf("UPDATE_OFFSET environment variable does not exist")
	}
	timeout, err := strconv.ParseInt(os.Getenv("UPDATE_TIMEOUT"), 10, 32)
	if err != nil {
		log.Println("UPDATE_TIMEOUT environment variable does not exist")
		return nil, fmt.Errorf("UPDATE_TIMEOUT environment variable does not exist")
	}
	upcfg := tgbotapi.UpdateConfig{
		Offset:  int(offset),
		Timeout: int(timeout),
	}

	return &GifBot{
		bot:          bot,
		updateConfig: upcfg,
		keyboards:    initKeyboards(),
		gfycat:       new(gfycat.Client),
	}, nil
}

func (b *GifBot) Start() {
	updates, err := b.bot.GetUpdatesChan(b.updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	prevCmd := ""
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msgText := ""
		keyboardNum := 0
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgText)

		switch update.Message.Text {
		case "/start":
			msgText = "Hi, I'm a simple gif bot, I can do simple bot things"
			prevCmd = ""
		case "Cats":
			msgText = "Do you want a meme?"
			keyboardNum = 1
			prevCmd = "Cats"
		case "Dogs":
			msgText = "Do you want a meme?"
			keyboardNum = 2
			prevCmd = "Dogs"
		case "Capybaras":
			msgText = "Do you want a meme?"
			keyboardNum = 3
			prevCmd = "Capybaras"
		case "Your Request":
			msgText = "Please, enter the request"
			keyboardNum = 4
			prevCmd = "Your Request"
		case "<- Back":
			msgText = "What do you want?"
			keyboardNum = 0
			prevCmd = ""
		case "Just Cats":
			if prevCmd == "Cats" {
				msgText = "Ok\n"
				keyboardNum = 1
				prevCmd = "Just Cats"
				msgText += getGifs("cats", b)
			}
		case "Just Dogs":
			if prevCmd == "Dogs" {
				msgText = "Ok\n"
				keyboardNum = 2
				prevCmd = "Just Dogs"
				msgText += getGifs("dogs", b)
			}
		case "Just Capybaras":
			if prevCmd == "Capybaras" {
				msgText = "Ok\n"
				keyboardNum = 3
				prevCmd = "Just Capybaras"
				msgText += getGifs("capybaras", b)
			}
		case "Cat memes":
			if prevCmd == "Cats" {
				msgText = "Ok\n"
				keyboardNum = 1
				prevCmd = "Cat memes"
				msgText += getGifs("cats meme", b)
			}
		case "Dog memes":
			if prevCmd == "Dogs" {
				msgText = "Ok\n"
				keyboardNum = 2
				prevCmd = "Dog memes"
				msgText += getGifs("dogs meme", b)
			}
		case "Capybaras memes":
			if prevCmd == "Capybaras" {
				msgText = "Ok\n"
				keyboardNum = 3
				prevCmd = "Capybaras memes"
				msgText += getGifs("capybaras meme", b)
			}
		default:
			if prevCmd == "Your Request" {
				request := update.Message.Text
				msgText += getGifs(request, b)
			} else if prevCmd == "I'm lucky" {
				request := "random"
				msgText += getGifs(request, b)
			} else {
				msgText = "You have send me something that I apparently can't understand."
				prevCmd = ""
			}
		}

		msg.ReplyMarkup = b.keyboards[keyboardNum]
		msg.Text = msgText
		if _, err := b.bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func getGifs(query string, b *GifBot) string {
	msgText := ""
	res, err := b.gfycat.Search(context.Background(), &gfycat.Request{
		Query:  query,
		InTest: true,

		MaxPageNumber: 1,
	})
	if err != nil {
		log.Fatal(err)
	}

	for page := range res.Pages {
		if page.Err != nil {
			fmt.Printf("%d: err: %v\n", page.PageNumber, page.Err)
			continue
		}

		for i, gfy := range page.Gfys {
			msgText += fmt.Sprintf("\t(#%d): %#v\n", i, gfy)
			log.Printf("\t(#%d): %#v\n", i, gfy)
		}
	}
	return msgText
}
