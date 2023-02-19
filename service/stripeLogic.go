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
	"github.com/taiki-nd/taxi_log/db"
	"github.com/taiki-nd/taxi_log/model"
)

type Product struct {
	ProductID string
	PriceID   string
	Price     int64
}

/**
 * GetUserInfoForStripe
 */
func GetUserInfoForStripe(uuid string) (string, error) {
	// user情報取得
	var user model.User
	err := db.DB.Table("users").Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return "", err
	}

	return user.StripeSubId, nil
}

/**
 * CreateSubscription
 */
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

/*
 * UpdateUserForStartSubscription
 */
func UpdateUserForStartSubscription(uuid string, subscription *stripe.Subscription) error {
	// user情報取得
	var user model.User
	err := db.DB.Table("users").Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return err
	}

	// user情報の更新
	err = db.DB.Model(&user).Updates(model.User{StripeCId: subscription.Customer.ID, StripeSubId: subscription.ID}).Error
	if err != nil {
		return err
	}

	return nil
}

/**
 * CancelSubscription
 */
func CancelSubscription(c *fiber.Ctx) (*stripe.Subscription, error) {
	fmt.Println("start cancel sub")
	// headerの確認
	var header AuthUser
	err := c.ReqHeaderParser(&header)
	if err != nil {
		log.Println("reqHeader parse error")
		return nil, err
	}
	uuid := header.Uuid

	// user情報取得
	var user model.User
	err = db.DB.Table("users").Where("uuid = ?", uuid).First(&user).Error
	if err != nil {
		return nil, err
	}

	// DBからcustomerIdとsubscriptionIdの取得
	stripeInfo := struct {
		StripeCId   string
		StripeSubId string
	}{}
	err = db.DB.Table("users").Where("uuid = ?", uuid).Find(&stripeInfo).Error
	if err != nil {
		return nil, fmt.Errorf("db_error")
	}

	if stripeInfo.StripeSubId == "" || stripeInfo.StripeCId == "" {
		return nil, fmt.Errorf("not_premium_plan")
	}

	stripe.Key = config.Config.StripeSecretKey

	// サブスクリプションのキャンセル
	s, err := sub.Cancel(stripeInfo.StripeSubId, nil)
	if err != nil {
		return nil, err
	}

	// カスタマーの削除
	cus, err := customer.Del(user.StripeCId, nil)
	if err != nil {
		return nil, err
	}
	log.Printf("success delete customer form stripe: %v", cus)

	err = db.DB.Debug().Model(&user).Updates(model.User{StripeCId: "", StripeSubId: ""}).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}
