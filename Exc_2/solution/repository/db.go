package repository

import (
	"ordersystem/model"
	"time"
)

type DatabaseHandler struct {
	// drinks represent all available drinks
	drinks []model.Drink
	// orders serves as order history
	orders []model.Order
}

// todo
func NewDatabaseHandler() *DatabaseHandler {
	// Init the drinks slice with some test data
	drinks := []model.Drink{
		{ID: 1, Name: "Coca-Cola", Price: 2.50, Description: "Classic soft drink"},
		{ID: 2, Name: "Pepsi", Price: 2.40, Description: "Popular cola beverage"},
		{ID: 3, Name: "Fanta", Price: 2.20, Description: "Orange flavoured drink"},
	}

	// Init orders slice with some test data (PRE-ORDERS)
	// Feel free to change the amounts/timestamps — these are just examples
	now := time.Now()
	orders := []model.Order{
		{ID: 1, DrinkID: 1, Amount: 10, CreatedAt: now.Add(-48 * time.Hour)}, // 10x Coca-Cola
		{ID: 2, DrinkID: 2, Amount: 7, CreatedAt: now.Add(-24 * time.Hour)},  // 7x Pepsi
		{ID: 3, DrinkID: 1, Amount: 5, CreatedAt: now.Add(-12 * time.Hour)},  // 5x Coca-Cola
		{ID: 4, DrinkID: 3, Amount: 3, CreatedAt: now.Add(-6 * time.Hour)},   // 3x Fanta
	}

	return &DatabaseHandler{
		drinks: drinks,
		orders: orders,
	}
}

// If your REST layer calls GetMenu(), keep this alias so it compiles.
func (db *DatabaseHandler) GetMenu() []model.Drink {
	return db.drinks
}

func (db *DatabaseHandler) GetDrinks() []model.Drink {
	return db.drinks
}

func (db *DatabaseHandler) GetOrders() []model.Order {
	return db.orders
}

// todo
func (db *DatabaseHandler) GetTotalledOrders() map[uint64]uint64 {
	// calculate total orders
	// key = DrinkID, value = Amount of orders
	totalledOrders := make(map[uint64]uint64)
	for _, o := range db.orders {
		totalledOrders[o.DrinkID] += uint64(o.Amount)
	}
	return totalledOrders
}

func (db *DatabaseHandler) AddOrder(order model.Order) {
	// todo
	// add order to db.orders slice
	if order.ID == 0 {
		order.ID = uint64(len(db.orders) + 1)
	}
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	db.orders = append(db.orders, order)
}
