package main

import (
	"encoding/binary"
	"fideliy/dins"
	"fmt"
	"time"
	tg "github.com/acteek/telegram-bot-api"
)

type Publisher struct {
	api   *dins.DinsApi
	bot   *tg.BotAPI
	store *Store
	Ch    chan Subscription
}

type Action int

const (
	Create Action = iota
	Delete
)

type Subscription struct {
	ChatID int64
	Action Action
}

//NewPublisher create new instance
func NewPublisher(api *dins.DinsApi, bot *tg.BotAPI, store *Store) *Publisher {
	return &Publisher{
		api:   api,
		bot:   bot,
		store: store,
		Ch:    make(chan Subscription),
	}
}

// Start Publisher and init subscribtions
func (p *Publisher) Start() {
	var tasks = make(map[int64](chan string))

	go func() {
		for {
			event := <-p.Ch
			switch event.Action {
			case Create:
				done := make(chan string)
				tasks[event.ChatID] = done
				p.startTask(event.ChatID, done)
			case Delete:
				tasks[event.ChatID] <- "done"
				close(tasks[event.ChatID])
				delete(tasks, event.ChatID)
			}
		}
	}()

	go p.initFromStore()
}

func (p *Publisher) startTask(chatID int64, done chan string) {
	ticker := time.NewTicker(5 * time.Second)
	fmt.Println("Task Create chatID: ", chatID)
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("Task Cancel chatID: ", chatID)
				return
			case <-ticker.C:
				user, _ p.store.Get(chatID)
				if  user.Subscriptions
				msg := tg.NewMessage(chatID, "tick")
				p.sendReply(msg)
			}
		}
	}()
}

func (p *Publisher) initFromStore() {
	for key := range p.store.Keys() {
		if len(key) == 0 {
			break
		}
		chatID := int64(binary.LittleEndian.Uint64(key))
		p.Ch <- Subscription{
			ChatID: chatID,
			Action: Create,
		}
	}
}

func (p *Publisher) sendReply(reply ...tg.Chattable) {
	for _, answer := range reply {
		if _, err := p.bot.Send(answer); err != nil {
			fmt.Println("Failed send message to telegram ", err)
		}
	}
}
