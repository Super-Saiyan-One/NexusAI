package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexus-ai/constant"
	stripeDto "nexus-ai/dto"
	"nexus-ai/model"
	"nexus-ai/repository"
	"nexus-ai/service"
	"nexus-ai/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/refund"
	"github.com/stripe/stripe-go/v81/webhook"
)

type StripeController interface {
	GetPaymentRepo() repository.PaymentRepository
	GetQuotaRepo() repository.QuotaRepository
	CreateCheckoutSession(c *gin.Context)
	Webhook(c *gin.Context)
	GetPaymentStatus(c *gin.Context)
	CreateRefund(c *gin.Context)
}

type stripeController struct {
	paymentService service.PaymentService
}

func NewStripeController() StripeController {
	paymentService := service.NewPaymentService()
	return &stripeController{paymentService: paymentService}
}

func (sc *stripeController) GetPaymentRepo() repository.PaymentRepository {
	return repository.NewPaymentRepository(model.GetDB())
}

func (sc *stripeController) GetQuotaRepo() repository.QuotaRepository {
	return repository.NewQuotaRepository(model.GetDB())
}

func (sc *stripeController) CreateCheckoutSession(c *gin.Context) {
	var req stripeDto.StripeCheckoutSessionRequest
	if constant.StripeKey == "" {
		utils.CommonError(c, http.StatusInternalServerError, "Stripe config error, please contact admin", constant.ErrorTypePaymentPrefix+"_stripe_create_checkout_session")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_create_checkout_session")
		return
	}
	if req.PaymentMethod == "alipay" {
		req.PriceID = constant.AlipayPriceID
	} else if req.PaymentMethod == "card" {
		req.PriceID = constant.CardPriceID
	} else {
		utils.CommonError(c, http.StatusBadRequest, "Invalid payment method", constant.ErrorTypePaymentPrefix+"_stripe_create_checkout_session")
		return
	}
	if req.PriceID == "" {
		utils.CommonError(c, http.StatusBadRequest, "Price config error, please contact admin", constant.ErrorTypePaymentPrefix+"_stripe_create_checkout_session")
		return
	}
	expireTime := time.Now().Add(30 * time.Minute)
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(req.PriceID),
				Quantity: stripe.Int64(req.Amount),
			},
		},
		SuccessURL:         stripe.String(constant.StripeSuccessURL),
		CancelURL:          stripe.String(constant.StripeCancelURL),
		ClientReferenceID:  stripe.String(c.GetString(string(constant.UserIDKey))),
		PaymentMethodTypes: stripe.StringSlice([]string{req.PaymentMethod}),
		ExpiresAt:          stripe.Int64(expireTime.Unix()),
	}
	session, err := session.New(params)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, "Failed to create checkout session: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_create_checkout_session")
		return
	}
	saveReq := stripeDto.StripePaymentCreateCheckoutSessionRequest{
		PaymentID:          utils.GenerateRandomUUID(16),
		UserID:             c.GetString(string(constant.UserIDKey)),
		PaymentPlatform:    "stripe", //
		PaymentScene:       req.PaymentScene,
		PaymentMethod:      req.PaymentMethod,
		PaymentCurrency:    string(session.Currency),
		PaymentAmount:      float64(session.AmountTotal),
		PaymentType:        "checkout_session",
		PaymentOrderNo:     session.ID, // stripe checkout session id
		PaymentTitle:       fmt.Sprintf("Use %s to %s by Stripe", req.PaymentMethod, req.PaymentScene),
		PaymentDescription: "Nexus AI Subscription",
		ExpireTime:         expireTime,
	}
	err = sc.paymentService.CreateCheckoutSession(sc.GetPaymentRepo(), &saveReq)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, "Failed to create checkout session: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_create_checkout_session")
		return
	}
	utils.CommonSuccess(c, http.StatusOK, "Checkout session created successfully", constant.ErrorTypePaymentPrefix+"_stripe_create_checkout_session", gin.H{
		"url":        session.URL,
		"session_id": session.ID,
	})
}

func (sc *stripeController) Webhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := c.GetRawData()
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to read request body: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_webhook")
		return
	}
	endpointSecret := constant.StripeWebhookSecret
	event, err := webhook.ConstructEvent(payload, c.GetHeader("Stripe-Signature"), endpointSecret)
	if err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Failed to verify signature: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_webhook")
		return
	}
	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			utils.CommonError(c, http.StatusBadRequest, "Failed to unmarshal checkout session data: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_webhook")
			return
		}
		userID := session.ClientReferenceID
		amount := session.AmountTotal
		currency := string(session.Currency)
		utils.LogInfo(c.Request.Context(), fmt.Sprintf("Checkout session completed for user %s, amount %d, currency %s", userID, amount, currency))

	}
}

func (sc *stripeController) GetPaymentStatus(c *gin.Context) {
	var req struct {
		SessionID string `json:"session_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_get_payment_status")
		return
	}
	if req.SessionID == "" {
		utils.CommonError(c, http.StatusBadRequest, "Session ID is required", constant.ErrorTypePaymentPrefix+"_stripe_get_payment_status")
		return
	}
	session, err := session.Get(req.SessionID, nil)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, "Failed to get checkout session status: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_get_payment_status")
		return
	}
	if session.ExpiresAt != 0 && time.Now().Unix() > session.ExpiresAt {
		utils.CommonError(c, http.StatusBadRequest, "Checkout session expired", constant.ErrorTypePaymentPrefix+"_stripe_get_payment_status")
		return
	}
	utils.CommonSuccess(c, http.StatusOK, "Checkout session retrieved successfully", constant.ErrorTypePaymentPrefix+"_stripe_get_payment_status", gin.H{
		"status": session.Status,
	})
}

func (sc *stripeController) CreateRefund(c *gin.Context) {
	var req struct {
		PaymentIntent string `json:"payment_intent" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.CommonError(c, http.StatusBadRequest, "Invalid request data: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_create_refund")
		return
	}
	if req.PaymentIntent == "" {
		utils.CommonError(c, http.StatusBadRequest, "Payment intent is required", constant.ErrorTypePaymentPrefix+"_stripe_create_refund")
		return
	}
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(req.PaymentIntent),
	}
	refund, err := refund.New(params)
	if err != nil {
		utils.CommonError(c, http.StatusInternalServerError, "Failed to create refund: "+err.Error(), constant.ErrorTypePaymentPrefix+"_stripe_create_refund")
		return
	}
	utils.CommonSuccess(c, http.StatusOK, "Refund created successfully", constant.ErrorTypePaymentPrefix+"_stripe_create_refund", gin.H{
		"refund": refund,
	})
}
