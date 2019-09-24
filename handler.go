package main

import (
	"fideliy/dins"
	"fideliy/helpers"
	"fmt"
	tg "github.com/acteek/telegram-bot-api"
	"log"
	"strings"
)

type Handler struct {
	api     *dins.DinsApi
	bot     *tg.BotAPI
	store   *Store
	baskets map[int64][]string
}

func NewHandler(api *dins.DinsApi, bot *tg.BotAPI, store *Store) *Handler {
	return &Handler{
		api:     api,
		bot:     bot,
		store:   store,
		baskets: make(map[int64][]string),
	}
}

func (h *Handler) sendReply(reply ...tg.Chattable) {
	for _, answer := range reply {
		if _, err := h.bot.Send(answer); err != nil {
			fmt.Println("Failed send message to telegram ", err)
		}
	}
}

func (h *Handler) callbackReply(cb *tg.CallbackQuery, text string) {
	if _, err := h.bot.AnswerCallbackQuery(tg.NewCallbackWithAlert(cb.ID, text)); err != nil {
		log.Println("Failed send answer callback to telegram", err)
	}
}

func (h *Handler) HandleCommand(msg *tg.Message) {
	reply := tg.NewMessage(msg.Chat.ID, "")
	switch msg.Command() {
	case "set_token":
		m := strings.Split(msg.Text, " ")
		if len(m) == 2 {
			token := m[1]
			user, er := h.api.GetUser(token)
			if er != nil {
				reply.Text = "Что то пошло не так попробуй другой"
			} else {
				if err := h.store.Put(msg.Chat.ID, user); err != nil {
					reply.Text = "Что то пошло не так попробуй другой"
				}
				reply.Text = user.Name + ", добро пожаловать"
			}

		} else {
			reply.Text = "Используй команду: /set_token your-token"
		}
	case "start":
		reply.Text = "Привет. Для авторизации используй \n /set_token {mydins-auth cookie c my.dins.ru}"
		reply.ReplyMarkup = helpers.BuildMainKeyboard()
	default:
		reply.Text = "Я не знаю такой команды"
	}

	h.sendReply(reply)

}

func (h *Handler) HandleMessage(msg *tg.Message) {
	reply := tg.NewMessage(msg.Chat.ID, "")
	switch msg.Text {
	case "Меню":
		if user, getErr := h.store.Get(msg.Chat.ID); getErr == nil {
			menu := h.api.GetMenu(user)
			if len(menu) == 0 {
				reply.Text = "Сейчас меню не доступно, попробуй позже"
			} else {
				reply.Text = "Вооот"
				reply.ReplyMarkup = helpers.BuildMenuKeyBoard(menu)
			}

		} else {
			reply.Text = "Ты кто такой ...? Используй: /set_token your-token"
		}

	case "Мои заказы":
		if user, getErr := h.store.Get(msg.Chat.ID); getErr == nil {
			orders := h.api.GetOrders(user)
			if len(orders) == 0 {
				reply.Text = "Ты ничего не заказал"
			} else {
				var names []string
				mealStore := h.api.CurrentMeals()
				for _, ord := range orders {
					names = append(names, mealStore[ord.MealID].Name)
				}

				reply.Text = strings.Join(names, ", ")
				reply.ReplyMarkup = helpers.BuildCancelOrderKeyBoard(orders[0])
			}

		} else {
			reply.Text = "Ты кто такой ...? Используй: /set_token your-token"
		}
	default:
		reply.Text = "🙀😴"
	}

	h.sendReply(reply)

}

func (h *Handler) HandleCallback(callback *tg.CallbackQuery) {
	switch data := callback.Data; {
	case data == "make_order":
		if basket, nonEmpty := h.baskets[callback.Message.Chat.ID]; nonEmpty {
			var names []string
			mealStore := h.api.CurrentMeals()
			for _, id := range basket {
				names = append(names, mealStore[id].Name)
			}
			submit := tg.NewMessage(callback.Message.Chat.ID, strings.Join(names, ", "))
			deleteMenu := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
			submit.ReplyMarkup = helpers.BuildOrderKeyBoard()
			h.sendReply(submit, deleteMenu)

		} else {
			h.callbackReply(callback, "Ты ничего не выбрал")
		}

	case data == "send_order":
		if basket, nonEmpty := h.baskets[callback.Message.Chat.ID]; nonEmpty {
			reply := tg.NewMessage(callback.Message.Chat.ID, "")
			user, getErr := h.store.Get(callback.Message.Chat.ID)
			if getErr != nil {
				reply.Text = "Что-то пошло не так, попробуй /set_token"
			} else {
				delSubmit := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
				h.sendReply(delSubmit)
				if err := h.api.SendOrder(basket, user); err != nil {
					reply.Text = "Что-то пошло не так"
					reply.ReplyMarkup = helpers.DinsRedirectKeyBoard(dinsEndpoint, "Заказать на сайте")
				} else {
					reply.Text = "Заказал для тебя"
				}
			}

			delete(h.baskets, callback.Message.Chat.ID)
			h.sendReply(reply)

		} else {
			h.callbackReply(callback, "Ты ничего не выбрал")
		}

	case data == "clear_order":
		reply := tg.NewMessage(callback.Message.Chat.ID, "Штош ...")
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)

		delete(h.baskets, callback.Message.Chat.ID)
		h.sendReply(reply, del)

	case strings.Contains(data, "cancel_order"):
		reply := tg.NewMessage(callback.Message.Chat.ID, "Штош...")
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)

		orderID := strings.Split(data, ":")[1]
		user, _ := h.store.Get(callback.Message.Chat.ID)

		if err := h.api.CancelOrder(orderID, user); err != nil {
			reply.Text = " Что то пошло не так, попробуй позже"
			h.sendReply(reply)
		} else {
			h.sendReply(reply, del)
		}

	default:
		h.baskets[callback.Message.Chat.ID] =
			append(h.baskets[callback.Message.Chat.ID], callback.Data)

		h.callbackReply(callback, "Добавил в корзину")
	}

}
