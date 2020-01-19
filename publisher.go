package main

import (
	"fideliy/dins"
	// "fideliy/helpers"
	"fmt"

	tg "github.com/acteek/telegram-bot-api"
	// "log"
	// "strconv"
)

type Publisher struct {
	api   *dins.DinsApi
	bot   *tg.BotAPI
	store *Store
	Ch    chan string
	keys  chan []byte
}

func NewPublisher(api *dins.DinsApi, bot *tg.BotAPI, store *Store) *Publisher {
	return &Publisher{
		api:   api,
		bot:   bot,
		store: store,
		Ch:    make(chan string),
		keys:  store.Keys(),
	}
}

func (p *Publisher) Start() {

	go func() {
		for {
			msg1 := <-p.Ch
			fmt.Println("received1", msg1)
		}
	}()

	for elem := range p.store.Keys() {
		if len(elem) == 0 {
			break
		}
		p.Ch <- string(elem)
	}

}

// func (h *Handler) sendReply(reply ...tg.Chattable) {
// 	for _, answer := range reply {
// 		if _, err := h.bot.Send(answer); err != nil {
// 			fmt.Println("Failed send message to telegram ", err)
// 		}
// 	}
// }
