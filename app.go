package main

import (
	"fideliy/dins"
	"flag"
	telegram "github.com/acteek/telegram-bot-api"
	"log"
)

func main() {
	log.Println("Starting...")
	var confPath string
	flag.StringVar(&confPath, "conf", "./config.json", "config path")
	flag.Parse()

	conf := FromFile(confPath)
	log.Println("With config", conf.Json())

	users := NewStore(conf.Store.Path)
	defer users.Close()

	bot, err := telegram.NewBotAPIWithEndpoint(conf.TgToken, conf.TgEndpoint)
	dinsApi := dins.NewDinsApi(conf.DinsEndpoint)
	handler := NewHandler(dinsApi, bot, users)

	if err != nil {
		log.Panic("Failed connect to telegram ", err)
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
			go handler.HandleCommand(update.Message)
		case m.Message != nil:
			go handler.HandleMessage(update.Message)
		case m.CallbackQuery != nil:
			go handler.HandleCallback(update.CallbackQuery)
		}

	}

}
