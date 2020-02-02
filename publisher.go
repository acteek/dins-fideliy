package main

import (
	"fideliy/dins"
	tg "github.com/acteek/telegram-bot-api"
	"log"
	"time"
)

//Publisher create and delete publish tasks for subscriptions
type Publisher struct {
	api   *dins.API
	bot   *tg.BotAPI
	store *Store
	Ch    chan Subscription
}

//Action for Subscripton
type Action int

// Subscription can be Create or Delete
const (
	Create Action = iota
	Delete
)

// Subscription data fields
type Subscription struct {
	ChatID int64
	Action Action
}

//NewPublisher create new instance
func NewPublisher(api *dins.API, bot *tg.BotAPI, store *Store) *Publisher {
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
				if _, isCreated := tasks[event.ChatID]; !isCreated {
					done := make(chan string)
					tasks[event.ChatID] = done
					p.subscriptionTask(event.ChatID, done)
				} else {
					log.Println("Task already created cahtID ", event.ChatID)
				}

			case Delete:
				tasks[event.ChatID] <- "done"
				close(tasks[event.ChatID])
				delete(tasks, event.ChatID)
			}
		}
	}()

	go p.initFromStore()
}


//TODO implement task logic
func (p *Publisher) subscriptionTask(chatID int64, done chan string) {
	ticker := time.NewTicker(5 * time.Second)
	log.Println("Task Create chatID: ", chatID)
	go func() {
		for {
			select {
			case <-done:
				log.Println("Task Cancel chatID: ", chatID)
				return
			case <-ticker.C:

				msg := tg.NewMessage(chatID, "tick")
				p.sendReply(msg)
			}
		}
	}()
}

func (p *Publisher) initFromStore() {
	for chatID := range p.store.ChatIDs() {
		user, _ := p.store.Get(chatID)
		if len(user.Subs) > 0 {
			p.Ch <- Subscription{
				ChatID: chatID,
				Action: Create,
			}
		}
	}
}

func (p *Publisher) sendReply(reply ...tg.Chattable) {
	for _, answer := range reply {
		if _, err := p.bot.Send(answer); err != nil {
			log.Println("Failed send message to telegram ", err)
		}
	}
}
