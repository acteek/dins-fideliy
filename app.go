package main

import (
	"fideliy/dins"
	"fideliy/helpers"
	telegram "github.com/acteek/telegram-bot-api"
	"log"
	"strings"
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
	baskets := make(map[int64][]string)

	defer users.Close()

	bot, err := telegram.NewBotAPI(botToken, tgEndpoint)
	dinsApi := dins.NewDinsApi(dinsEndpoint)
	if err != nil {
		log.Panic("Failed connect to telegram", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	//TODO only for tests
	users.Put(149199925, dins.User{ID: "1092", Name: "Sergey Ryazanov"})

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	// TODO   migrate logic to Handler
	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				msg := telegram.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Command() {
				case "set_token":
					m := strings.Split(update.Message.Text, " ")
					if len(m) == 2 {
						token := m[1]
						user, er := dinsApi.GetUser(token)
						log.Println(user)
						if er != nil {
							msg.Text = "Что то пошло не так попробуй другой"
						} else {
							if err := users.Put(update.Message.Chat.ID, user); err != nil {
								msg.Text = "Что то пошло не так попробуй другой"
							}
							msg.Text = user.Name + ", добро пожаловать"
						}

					} else {
						msg.Text = "Используй команду: /set_token your-token"
					}
				case "start":
					msg.Text = "Привет. Для авторизации используй \n /set_token {mydins-auth cookie c my.dins.ru}"
					msg.ReplyMarkup = helpers.BuildMainKeyboard()
				default:
					msg.Text = "Я не знаю такой команды"
				}
				if _, err := bot.Send(msg); err != nil {
					log.Panic(err)
				}

			} else {
				msg := telegram.NewMessage(update.Message.Chat.ID, "")
				switch update.Message.Text {
				case "Меню":
					if user, getErr := users.Get(update.Message.Chat.ID); getErr == nil {
						menu := dinsApi.GetMenu(user)
						if len(menu) == 0 {
							msg.Text = "Сейчас меню не доступно, попробуй позже"
						} else {
							msg.Text = "Вооот"
							msg.ReplyMarkup = helpers.BuildMenuKeyBoard(menu)
						}

					} else {
						msg.Text = "Ты кто такой ...? Используй: /set_token your-token"
					}

				case "Мои заказы":
					msg.Text = "Это пока не реализовано"
					msg.ReplyMarkup = helpers.DinsRedirectKeyBoard(dinsEndpoint, "Проверить на Сайте")
				default:
					msg.Text = "🙀😴"
				}

				if _, err := bot.Send(msg); err != nil {
					log.Panic("Failed Send message", err)
				}
			}
		} else if update.CallbackQuery != nil {

			switch update.CallbackQuery.Data {
			case "make_order":
				if basket, nonEmpty := baskets[update.CallbackQuery.Message.Chat.ID]; nonEmpty {
					var names []string
					mealStore := dinsApi.CurrentMeals()
					for _, id := range basket {
						names = append(names, mealStore[id].Name)
					}

					submit := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, strings.Join(names, ", "))
					submit.ReplyMarkup = helpers.BuildOrderKeyBoard()

					deleteMenu := telegram.NewDeleteMessage(
						update.CallbackQuery.Message.Chat.ID,
						update.CallbackQuery.Message.MessageID)

					if _, err := bot.Send(submit); err != nil {
						log.Panic("Failed Send message", err)
					}
					if _, err := bot.Send(deleteMenu); err != nil {
						log.Panic("Failed Send message", err)
					}

				} else {
					if _, err := bot.AnswerCallbackQuery(telegram.NewCallbackWithAlert(update.CallbackQuery.ID, "Ты ничего не выбрал")); err != nil {
						log.Panic("Failed Send message", err)
					}
				}

			case "send_order":
				if basket, nonEmpty := baskets[update.CallbackQuery.Message.Chat.ID]; nonEmpty {
					msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
					user, getErr := users.Get(update.CallbackQuery.Message.Chat.ID)
					if getErr != nil {
						msg.Text = "Что-то пошло не так, попробуй /set_token"
					} else {
						delSubmit := telegram.NewDeleteMessage(
							update.CallbackQuery.Message.Chat.ID,
							update.CallbackQuery.Message.MessageID)

						if err := dinsApi.SendOrder(basket, user); err != nil {
							msg.Text = "Что-то пошло не так"
							msg.ReplyMarkup = helpers.DinsRedirectKeyBoard(dinsEndpoint, "Заказать на сайте")
						} else {
							msg.Text = "Заказал для тебя"
						}
						if _, err := bot.Send(delSubmit); err != nil {
							log.Panic("Failed Send message", err)
						}
					}

					delete(baskets, update.CallbackQuery.Message.Chat.ID)

					if _, err := bot.Send(msg); err != nil {
						log.Panic("Failed Send message", err)
					}

				} else {
					if _, err := bot.AnswerCallbackQuery(telegram.NewCallbackWithAlert(update.CallbackQuery.ID, "Ты ничего не выбрал")); err != nil {
						log.Panic("Failed Send message", err)
					}
				}

			case "clear_order":
				msg := telegram.NewMessage(update.CallbackQuery.Message.Chat.ID, "Штош ...")
				sss := telegram.NewDeleteMessage(
					update.CallbackQuery.Message.Chat.ID,
					update.CallbackQuery.Message.MessageID)

				delete(baskets, update.CallbackQuery.Message.Chat.ID)

				if _, err := bot.Send(msg); err != nil {
					log.Panic("Failed Send message", err)
				}
				if _, err := bot.Send(sss); err != nil {
					log.Panic("Failed Send message", err)
				}

			default:
				baskets[update.CallbackQuery.Message.Chat.ID] =
					append(baskets[update.CallbackQuery.Message.Chat.ID], update.CallbackQuery.Data)

				if _, err := bot.AnswerCallbackQuery(telegram.NewCallbackWithAlert(update.CallbackQuery.ID, "Добавил в корзину")); err != nil {
					log.Panic("Failed Send message", err)
				}
			}

		}

	}

}
