package main

import (
	"fideliy/dins"
	"fideliy/helpers"
	"log"
	"time"

	tg "github.com/acteek/telegram-bot-api"
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
					log.Println("Task already created ChatID ", event.ChatID)
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

func (p *Publisher) subscriptionTask(chatID int64, done chan string) {
	ticker := time.NewTicker(time.Hour)
	log.Println("Task Create chatID: ", chatID)
	go func() {
		for {
			select {
			case <-done:
				log.Println("Task Cancel chatID: ", chatID)
				return
			case now := <-ticker.C:
				user, _ := p.store.Get(chatID)

				var active []string
				for sub, trigered := range user.Subs {
					tgHour := trigered.Truncate(time.Hour)
					if now.Truncate(time.Hour).After(tgHour) {
						active = append(active, sub)
					}
				}

				if contains(active, "Все Меню") {
					menu, _ := p.api.GetMenu(user)
					if len(menu) > 0 {
						msg := tg.NewMessage(chatID, "Меню по подписке")
						msg.ReplyMarkup = helpers.BuildMenuKeyBoard(menu)
						p.sendReply(msg)

						user.Subs["Все Меню"] = now
						p.store.Put(chatID, user)
					}

				} else if len(active) != 0 {
					menu, _ := p.api.GetMenu(user)
					var matches []string
					for _, meal := range menu {
						if contains(active, meal.Name) {
							matches = append(matches, meal.Name)
						}

						if len(matches) > 0 {
							msg := tg.NewMessage(chatID, "Меню по подписке "+matches[0])
							msg.ReplyMarkup = helpers.BuildMenuKeyBoard(menu)
							p.sendReply(msg)

							for _, sub := range matches {
								user.Subs[sub] = now

							}

							p.store.Put(chatID, user)
						}
					}

				}

			}
		}
	}()
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
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
