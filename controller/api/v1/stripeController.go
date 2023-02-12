package controller

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/taiki-nd/taxi_log/config"
	"github.com/taiki-nd/taxi_log/service"
	"github.com/taiki-nd/taxi_log/utils/constants"
	"google.golang.org/api/option"
)

func CreateSubscription(c *fiber.Ctx) error {
	// user認証
	statuses, errs, err := service.UserAuth(c)
	if err != nil {
		log.Printf("user auth error: %v", err)
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("user auth error: %v", err))
	}
	if len(errs) != 0 {
		log.Println(errs)
	}
	// signin確認
	if !statuses[0] {
		return service.ErrorResponse(c, []string{constants.USER_NOT_SIGININ}, "user not signin")
	}

	// headerの確認
	var header service.AuthUser
	err = c.ReqHeaderParser(&header)
	if err != nil {
		log.Println("reqHeader parse error")
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("header parse: %v", err))
	}
	//user_id := header.Id
	uuid := header.Uuid

	// ログイン情報の取得
	ctx := context.Background()
	opt := option.WithCredentialsFile(config.Config.FirebaseAuthPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Get an auth client from the firebase.App
	client, err := app.Auth(ctx)
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("firebase auth client error: %v", err))
	}

	user, err := client.GetUser(ctx, uuid)
	if err != nil {
		log.Printf("error getting user %s: %v\n", uuid, err)
		return service.ErrorResponse(c, []string{constants.USER_AUTH_ERROR}, fmt.Sprintf("firebase getting user error: %v", err))
	}
	email := user.Email

	// サブスクリプション処理
	sb, err := service.CreateSubscription(c, email, uuid)
	if err != nil {
		return service.ErrorResponse(c, []string{"create_subscription_err"}, fmt.Sprintf("create subscription err: %v", err))
	}

	// db登録
	err := 

	return service.SuccessResponse(c, []string{"success create subscription"}, sb, nil)

}

func CancelSubscription(c *fiber.Ctx) error {
	sub_id := c.Query("sub_id")
	stripe.Key = config.Config.StripeSecretKey
	s, _ := sub.Cancel(sub_id, nil)
	return service.SuccessResponse(c, []string{"success_cancel_subscription"}, s, nil)
}
