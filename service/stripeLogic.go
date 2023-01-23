package service

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	stripe "github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/price"
	"github.com/stripe/stripe-go/v74/setupintent"
	"github.com/stripe/stripe-go/v74/subscription"
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

/*
 * PaymentValidation
 */
func PaymentValidation(payment *model.Payment) (bool, []string) {
	var errs []string

	// PaymentMethodTypes
	if len(payment.PaymentMethodTypes) == 0 {
		errs = append(errs, "paymentMethodType_null_error")
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

func SetUpNewCard(c *fiber.Ctx, payment model.Payment) (interface{}, error) {
	// stripe接続
	stripe.Key = config.Config.StripeSecretKey

	var customerId *string
	c.ReqHeaderParser(&customerId)

	params := &stripe.SetupIntentParams{
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
		Customer: customerId,
	}
	pi, err := setupintent.New(params)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

func CancelSubscription(c *fiber.Ctx) interface{} {
	sub_id := c.Query("subscription_id")

	// stripe接続
	stripe.Key = config.Config.StripeSecretKey

	s, _ := subscription.Cancel(
		sub_id,
		nil,
	)
	fmt.Println(s)
	return s
}
