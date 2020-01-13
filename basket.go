package main

import (
	"fideliy/dins"
	"sync"
)

type chatID = int64
type itemID = string
type order = map[itemID]dins.Meal

//Basket is a struct wraper for basket map
type Basket struct {
	mx   sync.RWMutex
	data map[chatID]order
}


// NewBasket init  new empty basket
func NewBasket() *Basket {
	return &Basket{data: make(map[chatID]order)}

}

//Add meal to basket for chatID
func (b *Basket) Add(chatID chatID, meal dins.Meal) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if prev, has := b.data[chatID]; has {
		meal.Qty += prev[meal.ID].Qty
		prev[meal.ID] = meal
	} else {
		b.data[chatID] = order{meal.ID: meal}
	}

}

//Get order form basket for chatID
func (b *Basket) Get(chatID chatID) ([]dins.Meal, bool) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	order, has := b.data[chatID]
	var meals []dins.Meal
	for _, meal := range order {
		if meal.Qty > 0 {
			meals = append(meals, meal)
		}
	}

	return meals, has
}

//Delete order form basket for chatID
func (b *Basket) Delete(chatID chatID) {
	b.mx.Lock()
	defer b.mx.Unlock()
	delete(b.data, chatID)
}
