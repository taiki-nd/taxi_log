package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/taiki-nd/taxi_log/model"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
)

/**
 * GetProducts
 */
func GetProducts(c *fiber.Ctx) error {

	products, err := service.GetProducts()
	if err != nil {
		service.ErrorResponse(c, []string{"get_products_error"}, "failed to get products from stripe")
	}

	return service.SuccessResponse(c, nil, products, nil)
}

/**
 * CreateCustomer
 */
func CreateCustomer(c *fiber.Ctx) error {
	var customer *model.Customer

	// ボディーのパース
	err := c.BodyParser(&customer)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.BODY_PARSE_ERROR}, fmt.Sprintf("body parse error: %v", err))
	}

	// バリデーション
	_, errs := service.CustomerValidation(customer)
	if len(errs) != 0 {
		return service.ErrorResponse(c, errs, fmt.Sprintf("validation error: %v", errs))
	}

	// customer作成処理
	customerResponse, err := service.CreateCustomer(*customer)
	if err != nil {
		return service.ErrorResponse(c, []string{"create_customer_in_stripe_error"}, fmt.Sprintf("create customer error: %v", err))
	}

	return service.SuccessResponse(c, nil, customerResponse, nil)
}

/*
 * SetupPayMethod
 */
func SetupPayMethod(c *fiber.Ctx) error {
	var payment *model.Payment
	// ボディーのパース
	err := c.BodyParser(&payment)
	if err != nil {
		return service.ErrorResponse(c, []string{constants.BODY_PARSE_ERROR}, fmt.Sprintf("body parse error: %v", err))
	}

	// 支払い方法のセットアップ
	pi, err := service.SetUpNewCard(c, *payment)
	if err != nil {
		return service.ErrorResponse(c, []string{"setup_payment_method_in_stripe_error"}, fmt.Sprintf("setup payment error: %v", err))
	}

	return service.SuccessResponse(c, nil, pi, nil)
}

/*
 * CancelSubscription
 */
func CancelSubscription(c *fiber.Ctx) error {
	sub, err := service.CancelSubscription(c)
	if err != nil {
		return service.ErrorResponse(c, []string{"cancel_subscription_error"}, fmt.Sprintf("cancel subscription error: %v", err))
	}

	return service.SuccessResponse(c, nil, sub, nil)
}
