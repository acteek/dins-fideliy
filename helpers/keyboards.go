package helpers

import (
	"fideliy/dins"
	telegram "github.com/acteek/telegram-bot-api"
)

func BuildMainKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.NewReplyKeyboard(
		telegram.NewKeyboardButtonRow(
			telegram.NewKeyboardButton("Меню"),
			telegram.NewKeyboardButton("Мои заказы"),
		))
}

func BuildMenuKeyBoard(meals []dins.Meal) telegram.InlineKeyboardMarkup {
	var keyboard [][]telegram.InlineKeyboardButton

	orderButton := telegram.NewInlineKeyboardRow(
		telegram.NewInlineKeyboardButtonData("Перейти в корзину", "make_order"))

	for i := 0; i < len(meals); i++ {
		row := telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(meals[i].Name, meals[i].ID))
		keyboard = append(keyboard, row)

	}

	return telegram.InlineKeyboardMarkup{
		InlineKeyboard: append(keyboard, orderButton),
	}
}

func BuildOrderKeyBoard() telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Ок", "send_order"),
			telegram.NewInlineKeyboardButtonData("Я Передумал", "clear_order")))
}

func DinsRedirectKeyBoard(dinsEndpoint string, text string) telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonURL(text, dinsEndpoint+"/?page=fidel")))

}

func BuildCancelOrderKeyBoard(order dins.Orders) telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Отменить", "cancel_order:"+order.OrderID)))

}
