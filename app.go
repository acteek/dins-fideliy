package main

import (
	"fideliy/dins"
	"fmt"
	telegram "github.com/acteek/telegram-bot-api"
	"log"
	"strings"
)

const (
	botToken     = "987354230:AAGoDDLMxwowUY_wbuz6UCdgtD33eQE_Q4o"
	tgEndpoint   = "http://157.230.184.220/bot%s/%s" //proxy to api.telegram.com
	dinsEndpoint = "https://my.dins.ru/"
)

func main() {
	log.Println("Starting...")

	tokens := make(map[int64]string)

	bot, err := telegram.NewBotAPI(botToken, tgEndpoint)
	dinsApi := dins.NewDinsApi(dinsEndpoint)

	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message != nil {

			if update.Message.IsCommand() {
				msg := telegram.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "save_token":
					m := strings.Split(update.Message.Text, " ")
					if len(m) == 2 {
						tokens[update.Message.Chat.ID] = m[1]
						msg.Text = "token has been saved"
					} else {
						msg.Text = "please use format: save_token {your-token}"
					}

				case "get_token":
					if token, con := tokens[update.Message.Chat.ID]; con {
						msg.Text = "your token is " + token
					} else {
						msg.Text = "you don't have token yet"
					}
				case "menu":
					msg.Text = fmt.Sprintf("%v", dinsApi.GetMenu())
				default:
					msg.Text = "I don't know that command"
				}
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}

			} else {
				msg := telegram.NewMessage(update.Message.Chat.ID, "Дороу !")
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}
			}
		}

	}

}
