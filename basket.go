package main

import (
	"fideliy/dins"
	"sync"
)

type ChatId = int64
type ItemId = string

type Basket struct {
	Actions
	mx   sync.RWMutex
	data map[ChatId][]dins.Meal
}

type Actions interface {
	Add(chatId ChatId, meal dins.Meal)
	Get(chatId ChatId) ([]dins.Meal, bool)
	Delete(chatId ChatId)
}

func NewBasket() *Basket {
	return &Basket{data: make(map[ChatId][]dins.Meal)}

}

//TODO implement increment Qty
func (b *Basket) Add(chatId ChatId, meal dins.Meal) {
	b.mx.Lock()
	defer b.mx.Unlock()

	meals := b.data[chatId]
	b.data[chatId] = append(meals, meal)
}

func (b *Basket) Get(chatId ChatId) ([]dins.Meal, bool) {
	b.mx.RLock()
	defer b.mx.RUnlock()
	meals, has := b.data[chatId]

	return meals, has
}

func (b *Basket) Delete(chatId ChatId) {
	b.mx.Lock()
	defer b.mx.Unlock()
	delete(b.data, chatId)
}
