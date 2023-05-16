package purchase

import (
	"context"
	"errors"
	"fmt"
	"time"

	coffeeco "coffeeco/internal"
	"coffeeco/internal/loyalty"
	"coffeeco/internal/payment"
	"coffeeco/internal/store"

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

type StoreService interface {
	GetStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID) (float32, error)
}

type Service struct {
	cardService  CardChargeService
	purchaseRepo Repository
	storeService StoreService
}

func NewService(cardService CardChargeService, purchaseRepo Repository, storeService StoreService) *Service {
	return &Service{
		cardService:  cardService,
		purchaseRepo: purchaseRepo,
		storeService: storeService,
	}
}

func (s Service) CompletePurchase(ctx context.Context, storeID uuid.UUID, purchase *Purchase, coffeeBuxCard *loyalty.CoffeeBux) error {
	err := purchase.validateAndEnrich()
	if err != nil {
		return err
	}

	err = s.calculateStoreSpecificDiscount(ctx, storeID, purchase)
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
		err := coffeeBuxCard.Pay(ctx, purchase.ProductsToPurchase)
		if err != nil {
			return fmt.Errorf("failed to charge loyalty card: %w", err)
		}

	}

	err = s.purchaseRepo.Store(ctx, *purchase)
	if err != nil {
		return errors.New("failed to store purchase")
	}

	// coffeeBuxCard is a pointer because a customer
	// is under no obligratio to present a loyalty card
	// and therefore it can be nil
	if coffeeBuxCard != nil {
		coffeeBuxCard.AddStamp()
	}

	return nil
}

func (s *Service) calculateStoreSpecificDiscount(ctx context.Context, storeID uuid.UUID, purchase *Purchase) error {
	discount, err := s.storeService.GetStoreSpecificDiscount(ctx, storeID)
	if err != nil && !errors.Is(err, store.ErrNoDiscount) {
		return fmt.Errorf("failed to get discount: %w", err)
	}

	purchasePrice := purchase.total
	if discount > 0 {
		purchasePrice = *purchasePrice.Multiply(int64(100 - discount))
	}

	return nil
}
