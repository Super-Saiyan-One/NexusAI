package service

import (
	"fmt"
	stripeDto "nexus-ai/dto"
	dto "nexus-ai/dto/model"
	"nexus-ai/repository"
	"nexus-ai/utils"
)

type PaymentService interface {
	CreateCheckoutSession(repo repository.PaymentRepository, req *stripeDto.StripePaymentCreateCheckoutSessionRequest) error
}

type paymentService struct{}

func NewPaymentService() PaymentService {
	return &paymentService{}
}

func (ps *paymentService) CreateCheckoutSession(repo repository.PaymentRepository, req *stripeDto.StripePaymentCreateCheckoutSessionRequest) error {
	// 将 StripePaymentCreateCheckoutSessionRequest 转换为 Payment
	payment := &dto.Payment{
		PaymentID:          req.PaymentID,
		UserID:             req.UserID,
		PaymentPlatform:    req.PaymentPlatform,
		PaymentScene:       req.PaymentScene,
		PaymentMethod:      req.PaymentMethod,
		PaymentCurrency:    req.PaymentCurrency,
		PaymentAmount:      req.PaymentAmount,
		PaymentStatus:      0, // 初始状态为未支付
		PaymentType:        req.PaymentType,
		PaymentOrderNo:     req.PaymentOrderNo,
		PaymentTitle:       req.PaymentTitle,
		PaymentDescription: req.PaymentDescription,
		ExpireTime:         utils.MySQLTime(req.ExpireTime),
	}

	err := repo.Create(payment)
	if err != nil {
		return fmt.Errorf("failed to create payment record")
	}

	return nil
}
