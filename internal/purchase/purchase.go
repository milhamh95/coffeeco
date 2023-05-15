package purchase

import (
	coffeeco "coffeeco/internal"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"
	"context"
	"errors"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"
)

type Purchase struct {
	id                 uuid.UUID
	Store              store.Store
	ProductsToPurchase []coffeeco.Product
	total              money.Money
	PaymentMeans       payment.Means
	timeOfPurchase     time.Time
	CardToken          *string
}

// purchase is a pointer.
// because this function updates values that are missing
func (p *Purchase) validateAndEnrich() error {
	if len(p.ProductsToPurchase) == 0 {
		return errors.New("purchase must consist of at least one product")
	}

	p.total = *money.New(0, "USD")

	for _, v := range p.ProductsToPurchase {
		newTotal, _ := p.total.Add(&v.BasePrice)
		p.total = *newTotal
	}

	if p.total.IsZero() {
		return errors.New("purchase should never be 0. please validate")
	}

	p.id = uuid.New()
	p.timeOfPurchase = time.Now()

	return nil
}

type CardChargeService interface {
	ChargeCard(ctx context.Context, amount money.Money, cardToken string) error
}

type Service struct {
	cardService  CardChargeService
	purchaseRepo Repository
}

func (s Service) CompletePurchase(ctx context.Context, purchase *Purchase) error {
	err := purchase.validateAndEnrich()
	if err != nil {
		return err
	}

	switch purchase.PaymentMeans {
	case payment.MEANS_CARD:
		err := s.cardService.ChargeCard(ctx, purchase.total, *purchase.CardToken)
		if err != nil {
			return errors.New("card charge failed, cancelling purhcase")
		}
	case payment.MEANS_CASH:
		// TO DO
		return errors.New("payment method is not supported yet")
	case payment.MEANS_COFFEEBUX:
		err := s.purchaseRepo.Store(ctx, *purchase)
		if err != nil {
			return errors.New("failed to store purchase")
		}
	}

	return nil
}
