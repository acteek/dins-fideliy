package main

import (
	"fideliy/dins"
	hp "fideliy/helpers"
	"log"
	"strconv"
	"strings"
	"time"

	tg "github.com/acteek/telegram-bot-api"
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
			log.Println("Failed send message to telegram ", err)
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
		reply.ReplyMarkup = hp.BuildMainKeyboard()
	default:
		reply.Text = "Я не знаю такой команды"
	}

	h.sendReply(reply)

}

//HandleMessage handles telegram messages
func (h *Handler) HandleMessage(msg *tg.Message) {
	reply := tg.NewMessage(msg.Chat.ID, "")
	if user, getErr := h.store.Get(msg.Chat.ID); getErr == nil {

		switch msg.Text {
		case "Меню":
			menu, hasOrder := h.api.GetMenu(user)
			if len(menu) == 0 {
				reply.Text = "Сейчас меню не доступно, попробуй позже"
			} else if hasOrder {
				reply.Text = "Ты уже сделал заказ, используй \"Мои заказы\""
			} else {
				reply.Text = "Вооот"
				reply.ReplyMarkup = hp.BuildMenuKeyBoard(menu)
			}

		case "Мои Заказы":
			orders := h.api.GetOrders(user)
			if len(orders) == 0 {
				reply.Text = "Ты ничего не заказал"
			} else {
				var views []string
				mealStore := h.api.CurrentMeals
				for _, ord := range orders {
					view := mealStore[ord.MealID].Name + " " + ord.Qty + "шт."
					views = append(views, view)
				}

				reply.Text = strings.Join(views, ", ")
				reply.ReplyMarkup = hp.BuildCancelOrderKeyBoard(orders[0])
			}

		case "Подписки":
			reply.Text = "Подписки позволяют получать меню автоматически. Можно подписаться на все меню или конкретное блюдо "
			reply.ReplyMarkup = hp.BuildSubMainKeyBoard()
		default:
			reply.Text = "🙀😴"
			reply.ReplyMarkup = hp.BuildMainKeyboard()
		}
	} else {
		reply.Text = "Ты кто такой ...? Используй: /set_token your-token"
	}

	h.sendReply(reply)

}

//HandleCallback handles callbacks from keyboards
func (h *Handler) HandleCallback(callback *tg.CallbackQuery) {
	switch data := callback.Data; {
	case data == hp.MakeOrder:
		if order, nonEmpty := h.basket.Get(callback.Message.Chat.ID); nonEmpty {
			var views []string
			for _, meal := range order {
				view := meal.Name + " " + strconv.Itoa(meal.Qty) + "шт."
				views = append(views, view)
			}
			submit := tg.NewMessage(callback.Message.Chat.ID, strings.Join(views, ", "))
			deleteMenu := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
			submit.ReplyMarkup = hp.BuildOrderKeyBoard()
			h.sendReply(submit, deleteMenu)

		} else {
			h.callbackReply(callback, "Ты ничего не выбрал")
		}

	case data == hp.SendOrder:
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
					reply.ReplyMarkup = hp.DinsRedirectKeyBoard(h.api.Endpoint, "Заказать на сайте")
				} else {
					reply.Text = "Заказал для тебя"
				}
			}

			h.basket.Delete(callback.Message.Chat.ID)
			h.sendReply(reply)

		} else {
			h.callbackReply(callback, "Ты ничего не выбрал")
		}

	case data == hp.ClearOrder:
		reply := tg.NewMessage(callback.Message.Chat.ID, "Штош ...")
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)

		h.basket.Delete(callback.Message.Chat.ID)
		h.sendReply(reply, del)

	case data == hp.CloseMenu:
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		h.basket.Delete(callback.Message.Chat.ID)
		h.sendReply(del)

	case data == hp.Close:
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		h.sendReply(del)

	case data == hp.MakeSubs:
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		reply := tg.NewMessage(callback.Message.Chat.ID, "Вооот")

		reply.ReplyMarkup = hp.BuildMakeSubKeyBoard()
		h.sendReply(del, reply)

	case data == hp.CancelSubs:
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		reply := tg.NewMessage(callback.Message.Chat.ID, "Вооот")

		user, _ := h.store.Get(callback.Message.Chat.ID)

		var subNames []string

		for name := range user.Subs {
			subNames = append(subNames, name)
		}

		reply.ReplyMarkup = hp.BuildCancelSubKeyBoard(subNames)
		h.sendReply(del, reply)

	case data == hp.MakeSubsAll:
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		reply := tg.NewMessage(callback.Message.Chat.ID, "Создана подписка на все меню")

		user, _ := h.store.Get(callback.Message.Chat.ID)
		user.Subs["Все Меню"] = time.Time{}
		h.store.Put(callback.Message.Chat.ID, user)

		//TODO send message to Publisher

		h.sendReply(del, reply)

	case data == hp.CancelSubsAll:
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		reply := tg.NewMessage(callback.Message.Chat.ID, "Отменена подписка на все меню")

		user, _ := h.store.Get(callback.Message.Chat.ID)
		user.Subs = map[string]time.Time{}
		h.store.Put(callback.Message.Chat.ID, user)

		//TODO send message to Publisher

		h.sendReply(del, reply)

	case data == hp.MakeSubsMenu:
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		reply := tg.NewMessage(callback.Message.Chat.ID, "")

		user, _ := h.store.Get(callback.Message.Chat.ID)
		menu, _ := h.api.GetMenu(user)

		if len(menu) == 0 {
			reply.Text = "Сейчас меню не доступно, попробуй позже"
		} else {
			reply.ReplyMarkup = hp.BuildMakeSubMenuKeyBoard(menu)
		}

		//TODO send message to Publisher

		h.sendReply(del, reply)

	case strings.Contains(data, hp.MakeSub):
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		reply := tg.NewMessage(callback.Message.Chat.ID, "")

		mealID := hp.ParseValue(data)
		meal := h.api.CurrentMeals[mealID]

		user, _ := h.store.Get(callback.Message.Chat.ID)
		user.Subs[meal.Name] = time.Time{}
		h.store.Put(callback.Message.Chat.ID, user)

		reply.Text = "Создана подписка на " + meal.Name
		//TODO send message to Publisher

		h.sendReply(del, reply)

	case strings.Contains(data, hp.CancelSub):
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
		reply := tg.NewMessage(callback.Message.Chat.ID, "")

		mealName := hp.ParseValue(data)
		user, _ := h.store.Get(callback.Message.Chat.ID)
		delete(user.Subs, mealName)
		h.store.Put(callback.Message.Chat.ID, user)

		reply.Text = "Отменена подписка на " + mealName
		//TODO send message to Publisher

		h.sendReply(del, reply)

	case strings.Contains(data, hp.CancelOrder):
		reply := tg.NewMessage(callback.Message.Chat.ID, "Штош...")
		del := tg.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)

		orderID := hp.ParseValue(data)
		user, _ := h.store.Get(callback.Message.Chat.ID)

		if err := h.api.CancelOrder(orderID, user); err != nil {
			reply.Text = "Что то пошло не так, попробуй позже"
			h.sendReply(reply)
		} else {
			h.sendReply(reply, del)
		}

	case strings.Contains(data, hp.Order):
		mealID := hp.ParseValue(data)
		meal := h.api.CurrentMeals[mealID]
		meal.Qty = 1 // I can order only one item per iteration
		h.basket.Add(callback.Message.Chat.ID, meal)
		h.callbackReply(callback, "Добавил в корзину")

	default:
		log.Println("Don't match callback comand: ", data)

	}

}
