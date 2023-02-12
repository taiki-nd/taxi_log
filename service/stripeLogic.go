package service

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	stripe "github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/taiki-nd/taxi_log/config"
)

type Product struct {
	ProductID string
	PriceID   string
	Price     int64
}

func CreateSubscription(c *fiber.Ctx, email string, uid string) (*stripe.Subscription, error) {
	number := c.Query("card_number")
	expMonth := c.Query("exp_month")
	expYear := c.Query("exp_year")
	cvc := c.Query("cvc")

	expYear = "20" + expYear

	stripe.Key = config.Config.StripeSecretKey

	// 商品取得
	products := make([]Product, 0)
	priceParams := &stripe.PriceListParams{}
	priceIterator := price.List(priceParams)
	for priceIterator.Next() {
		products = append(products, Product{
			ProductID: priceIterator.Price().Product.ID,
			PriceID:   priceIterator.Price().ID,
			Price:     priceIterator.Price().UnitAmount,
		})
	}
	priceId := products[0].PriceID
	fmt.Printf("PriceID %v \n", priceId)

	// 顧客作成
	paramsCustomer := &stripe.CustomerParams{
		Name:        &uid,
		Email:       &email,
		Description: nil,
	}
	customer, err := customer.New(paramsCustomer)
	if err != nil {
		log.Printf("create_customer_error: %v", err)
		return nil, err
	}
	customerId := customer.ID
	fmt.Printf("customerId %v \n", customerId)

	// 支払い方法の作成
	paramsPaymentMethod := &stripe.PaymentMethodParams{
		Card: &stripe.PaymentMethodCardParams{
			Number:   stripe.String(number),
			ExpMonth: stripe.String(expMonth),
			ExpYear:  stripe.String(expYear),
			CVC:      stripe.String(cvc),
		},
		Type: stripe.String("card"),
	}
	paymentMethod, err := paymentmethod.New(paramsPaymentMethod)
	if err != nil {
		log.Printf("pay_method_error: %v", err)
		return nil, err
	}
	paymentMethodId := paymentMethod.ID
	fmt.Printf("paymentMethodId %v \n", paymentMethodId)

	// 支払い方法と顧客の紐付け
	params := &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerId),
	}
	attachedPaymentMethod, err := paymentmethod.Attach(paymentMethodId, params)
	if err != nil {
		log.Printf("pay_method_attached_error: %v", err)
		return nil, err
	}
	paymentMethodId = attachedPaymentMethod.ID

	// subscriptionの作成
	subscriptionParams := &stripe.SubscriptionParams{
		Customer: &customer.ID,
		Items: []*stripe.SubscriptionItemsParams{
			{
				Plan: &priceId,
			},
		},
		TrialEnd:             nil,
		DefaultPaymentMethod: &paymentMethodId,
	}
	sb, err := sub.New(subscriptionParams)
	if err != nil {
		if err != nil {
			log.Printf("create_subscription_error: %v", err)
			return nil, err
		}
		return nil, err
	}

	return sb, nil
}
