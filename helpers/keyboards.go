package helpers

import (
	"fideliy/dins"
	telegram "github.com/acteek/telegram-bot-api"
)

// BuildMainKeyboard returns main telegram keyboard
func BuildMainKeyboard() telegram.ReplyKeyboardMarkup {
	return telegram.NewReplyKeyboard(
		telegram.NewKeyboardButtonRow(
			telegram.NewKeyboardButton("Подписки"),
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

// BuildMakeSubMenuKeyBoard returns keyboard for subs menu
func BuildMakeSubMenuKeyBoard(meals []dins.Meal) telegram.InlineKeyboardMarkup {
	var keyboard [][]telegram.InlineKeyboardButton

	serviceButtons := telegram.NewInlineKeyboardRow(
		telegram.NewInlineKeyboardButtonData("❌ Закрыть", Close),
	)

	for i := 0; i < len(meals); i++ {
		row := telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(meals[i].Name, MakeSub+meals[i].ID))
		keyboard = append(keyboard, row)

	}

	return telegram.InlineKeyboardMarkup{
		InlineKeyboard: append(keyboard, serviceButtons),
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

// BuildSubMainKeyBoard returns main keyboard for subscriptions
func BuildSubMainKeyBoard() telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("🔔 Подписаться", MakeSubs),
			telegram.NewInlineKeyboardButtonData("🔕 Отписаться", CancelSubs),
			telegram.NewInlineKeyboardButtonData("🧾 Список", SubsList),
		),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("❌ Отмена", Close)),
	)
}

// BuildMakeSubKeyBoard returns make keyboard for subscriptions
func BuildMakeSubKeyBoard() telegram.InlineKeyboardMarkup {
	return telegram.NewInlineKeyboardMarkup(
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("На всё меню", MakeSubsAll),
			telegram.NewInlineKeyboardButtonData("Выбрать блюда", MakeSubsMenu),
		),
		telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData("❌ Отмена", Close)),
	)
}

// BuildCancelSubKeyBoard returns make keyboard for subscriptions
func BuildCancelSubKeyBoard(subs []string) telegram.InlineKeyboardMarkup {
	var keyboard [][]telegram.InlineKeyboardButton

	serviceButton := telegram.NewInlineKeyboardRow(
		telegram.NewInlineKeyboardButtonData("❌ Закрыть", Close),
	)

	cancelAllButton := telegram.NewInlineKeyboardRow(
		telegram.NewInlineKeyboardButtonData("От всего", CancelSubsAll),
	)

	for i := 0; i < len(subs); i++ {
		row := telegram.NewInlineKeyboardRow(
			telegram.NewInlineKeyboardButtonData(subs[i], CancelSub+subs[i]))
		keyboard = append(keyboard, row)

	}

	return telegram.InlineKeyboardMarkup{
		InlineKeyboard: append(keyboard, cancelAllButton, serviceButton),
	}

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
