package service

import (
	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/model"
)

func GetProducts() ([]model.Product, error) {
	// stripe接続
	stripe.Key = config.Config.StripeSecretKey

	products := make([]model.Product, 0)

	priceParams := &stripe.PriceListParams{}
	priceIterator := price.List(priceParams)
	for priceIterator.Next() {
		products = append(products, model.Product{
			Id:      priceIterator.Price().Product.ID,
			PriceId: priceIterator.Price().ID,
			Price:   priceIterator.Price().UnitAmount,
		})
	}

	return products, nil
}
