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

var layout = "15:04"

var evening = TimeRange{
	Start: "16:00",
	End:   "22:00",
}

var morning = TimeRange{
	Start: "09:00",
	End:   "11:50",
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
			case nowUTC := <-ticker.C:
				user, _ := p.store.Get(chatID)

				loc, _ := time.LoadLocation("Europe/Moscow")
				now := nowUTC.In(loc)
				nowTime := now.Format(layout)

				var active []string
				for sub, trigered := range user.Subs {
					tgHour := trigered.In(loc).Truncate(24 * time.Hour)

					if now.Truncate(24 * time.Hour).After(tgHour) && evening.contain(nowTime)  {
						active = append(active, sub)
					}
				}

				if contains(active, "Все Меню") {
					if menu, hasOrder := p.api.GetMenu(user); len(menu) > 0 && !hasOrder {
						msg := tg.NewMessage(chatID, "Меню по подписке")
						msg.ReplyMarkup = helpers.BuildMenuKeyBoard(menu)
						p.sendReply(msg)

						user.Subs["Все Меню"] = now
						p.store.Put(chatID, user)

					}

				} else if len(active) != 0 {
					if menu, hasOrder := p.api.GetMenu(user); len(menu) > 0 && !hasOrder {

						var matches []string
						for _, meal := range menu {
							if contains(active, meal.Name) {
								matches = append(matches, meal.Name)
							}
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
