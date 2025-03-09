package router

import (
	"nexus-ai/controller"
	"nexus-ai/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRouter(server *gin.Engine) {
	stripeController := controller.NewStripeController()
	paymentRouter := server.Group("/payment")
	stripeRouter := paymentRouter.Group("/stripe")
	stripeRouter.Use(middleware.UserVerifyMiddleware())
	{
		stripeRouter.POST("/create_checkout_seesion", stripeController.CreateCheckoutSession)
		stripeRouter.POST("/webhook", stripeController.Webhook)
		stripeRouter.GET("/get_payment_status", stripeController.GetPaymentStatus)
		stripeRouter.POST("/create_refund", stripeController.CreateRefund)
	}
}
