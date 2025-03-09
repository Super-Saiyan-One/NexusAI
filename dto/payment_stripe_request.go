package dto

import "time"

type StripeCheckoutSessionRequest struct {
	PriceID       string `json:"price_id"`
	Amount        int64  `json:"amount" binding:"required"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	PaymentScene  string `json:"payment_scene" binding:"required"`
}

type StripePaymentCreateCheckoutSessionRequest struct {
	PaymentID          string    `json:"payment_id"`          // 支付记录唯一标识
	UserID             string    `json:"user_id"`             // 关联的用户ID
	PaymentPlatform    string    `json:"payment_platform"`    // 支付平台
	PaymentScene       string    `json:"payment_scene"`       // 支付场景
	PaymentMethod      string    `json:"payment_method"`      // 支付方式
	PaymentCurrency    string    `json:"payment_currency"`    // 支付币种
	PaymentAmount      float64   `json:"payment_amount"`      // 支付金额
	PaymentType        string    `json:"payment_type"`        // 支付类型
	PaymentOrderNo     string    `json:"payment_order_no"`    // 支付订单号
	PaymentTitle       string    `json:"payment_title"`       // 支付标题
	PaymentDescription string    `json:"payment_description"` // 支付描述
	ExpireTime         time.Time `json:"expire_time"`         // 支付过期时间
}
