package loyalty

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/store"

	"github.com/google/uuid"
)

var N_MINIMUM_PURCHASE_FOR_FREE_DRINK = 10
var N_FREE_DRINK_FOR_CUSTOMER = 1

type CoffeeBux struct {
	ID                                   uuid.UUID
	store                                store.Store
	coffeeLover                          coffeeco.CoffeeLover
	FreeDrinksAvailable                  int
	RemainingDrinkPurchaseUntilFreeDrink int
}

func (c *CoffeeBux) AddStamp() {
	if c.RemainingDrinkPurchaseUntilFreeDrink == N_FREE_DRINK_FOR_CUSTOMER {
		c.RemainingDrinkPurchaseUntilFreeDrink = N_MINIMUM_PURCHASE_FOR_FREE_DRINK
		c.FreeDrinksAvailable += N_FREE_DRINK_FOR_CUSTOMER
		return
	}

	c.RemainingDrinkPurchaseUntilFreeDrink--
}
