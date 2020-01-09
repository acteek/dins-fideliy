package main

import (
	"fideliy/dins"
	"sync"
)

type ChatId = int64
type ItemId = string
type Order = map[ItemId]dins.Meal

type Basket struct {
	Actions
	mx   sync.RWMutex
	data map[ChatId]Order
}

type Actions interface {
	Add(chatId ChatId, meal dins.Meal)
	Get(chatId ChatId) ([]dins.Meal, bool)
	Delete(chatId ChatId)
}

func NewBasket() *Basket {
	return &Basket{data: make(map[ChatId]Order)}

}

func (b *Basket) Add(chatId ChatId, meal dins.Meal) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if order, has := b.data[chatId]; has {
		meal.Qty += order[meal.ID].Qty
		order[meal.ID] = meal
	} else {
		b.data[chatId] = Order{meal.ID: meal}
	}

}

func (b *Basket) Get(chatId ChatId) ([]dins.Meal, bool) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	order, has := b.data[chatId]
	var meals []dins.Meal
	for _, meal := range order {
		if meal.Qty > 0 {
			meals = append(meals, meal)
		}
	}

	return meals, has
}

func (b *Basket) Delete(chatId ChatId) {
	b.mx.Lock()
	defer b.mx.Unlock()
	delete(b.data, chatId)
}
