package main

import (
	"fideliy/dins"
	telegram "github.com/acteek/telegram-bot-api"
	"log"
)

const (
	botToken = "880777116:AAEOXE6RVHEWzzC0hthVPjwG37WxsUbRy2U" // Prod Bot
	//botToken     = "987354230:AAGoDDLMxwowUY_wbuz6UCdgtD33eQE_Q4o" // Test_Bot
	tgEndpoint   = "http://157.230.184.220/bot%s/%s" //proxy to api.telegram.com
	dinsEndpoint = "https://my.dins.ru"
)

func main() {
	log.Println("Starting...")

	users := NewStore("./data")
	defer users.Close()

	bot, err := telegram.NewBotAPI(botToken, tgEndpoint)
	dinsApi := dins.NewDinsApi(dinsEndpoint)
	handler := NewHandler(dinsApi, bot, users)

	if err != nil {
		log.Panic("Failed connect to telegram", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic("Failed to get updates from telegram")
	}

	for update := range updates {
		switch m := update; {
		case m.Message != nil && update.Message.IsCommand():
			handler.HandleCommand(update.Message)
		case m.Message != nil:
			handler.HandleMessage(update.Message)
		case m.CallbackQuery != nil:
			handler.HandleCallback(update.CallbackQuery)
		}

	}

}
