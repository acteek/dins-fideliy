package helpers

import (
	"fideliy/dins"
	telegram "github.com/acteek/telegram-bot-api"
)

// BuildMainKeyboard returns main telegram keyboard
func BuildMainKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.NewReplyKeyboard(
		telegram.NewKeyboardButtonRow(
			// telegram.NewKeyboardButton("Подписки"),
			telegram.NewKeyboardButton("Мои Заказы"),
			telegram.NewKeyboardButton("Меню"),
		))
}

// BuildMenuKeyBoard returns keyboard for menu
func BuildMenuKeyBoard(meals []dins.Meal) telegram.InlineKeyboardMarkup {
	var keyboard [][]telegram.InlineKeyboardButton

	orderButton := telegram.NewInlineKeyboardRow(
		telegram.NewInlineKeyboardButtonData("❌ Отмена", CloseMenu),
		telegram.NewInlineKeyboardButtonData("✅ В корзину", MakeOrder),
	)

	for i := 0; i < len(meals); i++ {
		row := telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(meals[i].Name, Order+meals[i].ID))
		keyboard = append(keyboard, row)

	}

	return telegram.InlineKeyboardMarkup{
		InlineKeyboard: append(keyboard, orderButton),
	}
}

// BuildOrderKeyBoard returns keyboard for make order
func BuildOrderKeyBoard() telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("❌ Я Передумал", ClearOrder),
			telegram.NewInlineKeyboardButtonData("✅ Ок", SendOrder),
		))
}

// BuildSubKeyBoard returns main keyboard for subscriptions
func BuildSubKeyBoard() telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("🔔 Подписаться", MakeSubs),
			telegram.NewInlineKeyboardButtonData("🔕 Отписаться", CancelSubs),
		),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("❌ Отмена", Close)),
	)
}

// DinsRedirectKeyBoard returns keyboard with redirect to my.dins.ru
func DinsRedirectKeyBoard(dinsEndpoint string, text string) telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonURL(text, dinsEndpoint+"/?page=fidel")))

}

// BuildCancelOrderKeyBoard returns keyboard with orders
func BuildCancelOrderKeyBoard(order dins.Order) telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("Отменить", CancelOrder+order.ID)))

}
