package service

import (
	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/model"
)

func CustomerValidation(customer *model.Customer) (bool, []string) {
	var errs []string

	// name
	if len(customer.Name) == 0 {
		errs = append(errs, "name_null_error")
	}

	// email
	if len(customer.Email) == 0 {
		errs = append(errs, "email_null_error")
	}

	// description
	if len(customer.Description) == 0 {
		errs = append(errs, "description_null_error")
	}

	// errの出力
	if len(errs) != 0 {
		return false, errs
	}

	return true, errs
}

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

func CreateCustomer(c model.Customer) (interface{}, error) {

	// stripe接続
	stripe.Key = config.Config.StripeSecretKey

	params := &stripe.CustomerParams{
		Name:        &c.Name,
		Email:       &c.Email,
		Description: &c.Description,
	}
	customerFromStripe, err := customer.New(params)
	if err != nil {
		return nil, err
	}

	return &customerFromStripe, nil
}
