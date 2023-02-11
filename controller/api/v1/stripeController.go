package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/service"
)

func CreateSubscription(c *fiber.Ctx) error {
	// ログイン情報の取得

}

func CancelSubscription(c *fiber.Ctx) error {
	sub_id := c.Query("sub_id")
	stripe.Key = config.Config.StripeSecretKey
	s, _ := sub.Cancel(sub_id, nil)
	return service.SuccessResponse(c, []string{"success_cancel_subscription"}, s, nil)
}
