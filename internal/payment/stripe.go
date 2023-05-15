package payment

import (
	"errors"

	"github.com/stripe/stripe-go/v74/client"
)

type StripeService struct {
	stripeClient *client.API
}

func NewStripeService(apiKey string) (*StripeService, error) {
	if apiKey == "" {
		return nil, errors.New("API key cannot be nil")
	}
	sc := &client.API{}
	sc.Init(apiKey, nil)

	return &StripeService{stripeClient: sc}, nil
}
