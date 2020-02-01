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
	CancelSubs string = "cancel_subscription"

	//MakeSubs make subscription menu
	MakeSubs string = "make_subscription"
)

//ParseValue parse value for callback with value
func ParseValue(data string) string {
	return strings.Split(data, ":")[1]
}
