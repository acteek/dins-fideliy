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
	log.Println("With config", conf.JSON())

	users := NewStore(conf.Store.Path)
	defer users.Close()

	bot, err := telegram.NewBotAPIWithEndpoint(conf.TgToken, conf.TgEndpoint)
	dinsAPI := dins.NewDinsAPI(conf.DinsEndpoint)
	handler := NewHandler(dinsAPI, bot, users)
	publisher := NewPublisher(dinsAPI, bot, users)

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

	publisher.Start()

	for update := range updates {
		switch m := update; {
		case m.Message != nil && update.Message.IsCommand():
			go handler.HandleCommand(update.Message)
		case m.Message != nil:
			publisher.Ch <- m.Message.Text
			go handler.HandleMessage(update.Message)
		case m.CallbackQuery != nil:
			go handler.HandleCallback(update.CallbackQuery)
		}

	}

}
