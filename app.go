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
	if err != nil {
		log.Panic("Failed connect to telegram", err)
	}

	handler := NewHandler(dinsApi, bot, users)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//TODO only for tests
	users.Put(149199925, dins.User{ID: "1092", Name: "Sergey Ryazanov"})

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		switch m := update; {
		case m.Message != nil:
			if update.Message.IsCommand() {
				handler.HandleCommand(update.Message)
			} else {
				handler.HandleMessage(update.Message)
			}
		case m.CallbackQuery != nil:
			handler.HandleCallback(update.CallbackQuery)
		}

	}

}
