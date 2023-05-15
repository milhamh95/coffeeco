package loyalty

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/purchase"
	"coffeeco/internal/store"
	"context"
	"errors"
	"fmt"

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

func (c *CoffeeBux) Pay(ctx context.Context, purchases []purchase.Purchase) error {
	lp := len(purchases)
	if lp == 0 {
		return errors.New("nothing to buy")
	}

	if c.FreeDrinksAvailable < lp {
		return fmt.Errorf("not enough coffeeBux to cover entire purchase. Have %d, need %d", lp, c.FreeDrinksAvailable)
	}

	c.FreeDrinksAvailable = c.FreeDrinksAvailable - lp
	return nil
}
