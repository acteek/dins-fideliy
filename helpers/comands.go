package helpers

import "strings"

//Callabck comands from InlineKeyboards
const (

	//Close -  common Comand
	Close string = "close"

	//CloseMenu -  close meal menu
	CloseMenu string = "close_menu"

	//Order - user add order to basket
	Order = "order:"

	//MakeOrder - user ready to make order/check basket
	MakeOrder string = "make_order"

	//ClearOrder - clear basket and cancel order
	ClearOrder string = "clear_order"

	//SendOrder - send order to server
	SendOrder string = "send_order"

	//CancelOrder - cancel order on server
	CancelOrder string = "cancel_order:"

	//CancelSubs - cancel subscribtion menu
	CancelSubs string = "cancel_subscriptions"

	//MakeSubs make subscription menu
	MakeSubs string = "make_subscriptions"


	//list of  subscriptions 
	SubsList string = "subscriptions_list"

	//make subscription for menu
	MakeSubsAll string = "make_subscription_all"

	//cancel subscription for menu
	CancelSubsAll string = "cancel_subscription_all"

	// make subscription for particular meal
	MakeSubsMenu string = "make_subscriptions_menu"

	//MakeSub with meal ID
	MakeSub = "make_subscription:"

	//CancelSub with meal ID
	CancelSub = "cancel_subscription:"
)

//ParseValue parse value for callback with value
func ParseValue(data string) string {
	return strings.Split(data, ":")[1]
}
