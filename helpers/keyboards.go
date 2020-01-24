package helpers

import (
	"fideliy/dins"
	telegram "github.com/acteek/telegram-bot-api"
)

func BuildMainKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.NewReplyKeyboard(
		telegram.NewKeyboardButtonRow(
			telegram.NewKeyboardButton("Подписки"),
			telegram.NewKeyboardButton("Заказы"),
			telegram.NewKeyboardButton("Меню"),
		))
}

func BuildMenuKeyBoard(meals []dins.Meal) telegram.InlineKeyboardMarkup {
	var keyboard [][]telegram.InlineKeyboardButton

	orderButton := telegram.NewInlineKeyboardRow(
		telegram.NewInlineKeyboardButtonData("Отмена", "close_menu"),
		telegram.NewInlineKeyboardButtonData("В корзину", "make_order"),
	)

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
			telegram.NewInlineKeyboardButtonData("Я Передумал", "clear_order"),
			telegram.NewInlineKeyboardButtonData("Ок", "send_order"),
		))
}

func BuildSubKeyBoard() telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Подписаться", "make_subscription")),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Отписаться", "cancel_subscription")),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Отмена", "close")),
	)
}

func DinsRedirectKeyBoard(dinsEndpoint string, text string) telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonURL(text, dinsEndpoint+"/?page=fidel")))

}

func BuildCancelOrderKeyBoard(order dins.Order) telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Отменить", "cancel_order:"+order.ID)))

}
