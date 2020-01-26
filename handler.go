package main

import (
	"fideliy/dins"
	"fideliy/helpers"
	"fmt"
	tg "github.com/acteek/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

//Handler describes all methods for handle message from telegram
type Handler struct {
	api    *dins.API
	bot    *tg.BotAPI
	store  *Store
	basket *Basket
}

//NewHandler returns new Handler instance 
func NewHandler(api *dins.API, bot *tg.BotAPI, store *Store) *Handler {
	return &Handler{
		api:    api,
		bot:    bot,
		store:  store,
		basket: NewBasket(),
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
//HandleCommand handles telegram commands
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
				reply.Text = user.Name + ", добро пожаловать !"
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

//HandleMessage handles telegram messages
func (h *Handler) HandleMessage(msg *tg.Message) {
	reply := tg.NewMessage(msg.Chat.ID, "")
	switch msg.Text {
	case "Меню":
		if user, getErr := h.store.Get(msg.Chat.ID); getErr == nil {
			menu, hasOrder := h.api.GetMenu(user)
			if len(menu) == 0 {
				reply.Text = "Сейчас меню не доступно, попробуй позже"
			} else if hasOrder {
				reply.Text = "Ты уже сделал заказ, используй \"Мои заказы\""
			} else {
				reply.Text = "Вооот"
				reply.ReplyMarkup = helpers.BuildMenuKeyBoard(menu)
			}

		} else {
			reply.Text = "Ты кто такой ...? Используй: /set_token your-token"
		}

	case "Мои Заказы":
		if user, getErr := h.store.Get(msg.Chat.ID); getErr == nil {
			orders := h.api.GetOrders(user)
			if len(orders) == 0 {
				reply.Text = "Ты ничего не заказал"
			} else {
				var views []string
				mealStore := h.api.CurrentMeals
				for _, ord := range orders {
					view := mealStore[ord.MealID].Name + " X" + ord.Qty
					views = append(views, view)
				}

				reply.Text = strings.Join(views, ", ")
				reply.ReplyMarkup = helpers.BuildCancelOrderKeyBoard(orders[0])
			}

		} else {
			reply.Text = "Ты кто такой ...? Используй: /set_token your-token"
		}
	case "Подписки":
		if _, getErr := h.store.Get(msg.Chat.ID); getErr == nil {

			reply.Text = "Можно подписаться на все меню или конкретное блюдо "
			reply.ReplyMarkup = helpers.BuildSubKeyBoard()

		} else {
			reply.Text = "Ты кто такой ...? Используй: /set_token your-token"
		}
	default:
		reply.Text = "🙀😴"
		reply.ReplyMarkup = helpers.BuildMainKeyboard()
	}

	h.sendReply(reply)

}

//HandleCallback handles callbacks from keyboards
func (h *Handler) HandleCallback(callback *tg.CallbackQuery) {
	switch data := callback.Data; {
	case data == "make_order":
		if order, nonEmpty := h.basket.Get(callback.Message.Chat.ID); nonEmpty {
			var views []string
			for _, meal := range order {
				view := meal.Name + " X" + strconv.Itoa(meal.Qty)
				views = append(views, view)
			}
			submit := tg.NewMessage(callback.Message.Chat.ID, strings.Join(views, ", "))
			deleteMenu := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
			submit.ReplyMarkup = helpers.BuildOrderKeyBoard()
			h.sendReply(submit, deleteMenu)

		} else {
			h.callbackReply(callback, "Ты ничего не выбрал")
		}

	case data == "send_order":
		if order, nonEmpty := h.basket.Get(callback.Message.Chat.ID); nonEmpty {
			reply := tg.NewMessage(callback.Message.Chat.ID, "")
			user, getErr := h.store.Get(callback.Message.Chat.ID)
			if getErr != nil {
				reply.Text = "Что-то пошло не так, попробуй /set_token"
			} else {
				delSubmit := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
				h.sendReply(delSubmit)
				if err := h.api.SendOrder(order, user); err != nil {
					reply.Text = "Что-то пошло не так"
					reply.ReplyMarkup = helpers.DinsRedirectKeyBoard(h.api.Endpoint, "Заказать на сайте")
				} else {
					reply.Text = "Заказал для тебя"
				}
			}

			h.basket.Delete(callback.Message.Chat.ID)
			h.sendReply(reply)

		} else {
			h.callbackReply(callback, "Ты ничего не выбрал")
		}

	case data == "clear_order":
		reply := tg.NewMessage(callback.Message.Chat.ID, "Штош ...")
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)

		h.basket.Delete(callback.Message.Chat.ID)
		h.sendReply(reply, del)

	case data == "close_menu":
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		h.basket.Delete(callback.Message.Chat.ID)
		h.sendReply(del)

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

	// if didn't math any comands callback.Data contains itemId
	default:
		meal := h.api.CurrentMeals[callback.Data]
		meal.Qty = 1 // I can order only one item per iteration
		h.basket.Add(callback.Message.Chat.ID, meal)
		h.callbackReply(callback, "Добавил в корзину")
	}

}
