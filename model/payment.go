package model

import (
	"nexus-ai/common"
	"nexus-ai/utils"

	"gorm.io/gorm"
)

// 系统中所有支付交易记录，包括支付信息、账单信息和退款信息等
type Payment struct {
	PaymentID string `gorm:"column:payment_id;type:char(36);primaryKey;default:(UUID())" json:"payment_id"`      // 支付记录唯一标识
	UserID    string `gorm:"column:user_id;type:char(36);index;not null;foreignKey:User(UserID)" json:"user_id"` // 关联的用户ID

	PaymentPlatform      string  `gorm:"column:payment_platform;size:50;index;not null" json:"payment_platform"`               // 支付平台(stripe/paddle/alipay/wechat等)
	PaymentScene         string  `gorm:"column:payment_scene;size:20;index;not null" json:"payment_scene"`                     // 支付场景(充值/订阅等)
	PaymentMethod        string  `gorm:"column:payment_method;size:50;not null" json:"payment_method"`                         // 支付方式(信用卡/银行转账/支付宝/微信等)
	PaymentCurrency      string  `gorm:"column:payment_currency;size:10;not null" json:"payment_currency"`                     // 支付币种
	PaymentAmount        float64 `gorm:"column:payment_amount;type:decimal(10,2);not null;default:0.00" json:"payment_amount"` // 支付金额
	PaymentStatus        int8    `gorm:"column:payment_status;index;not null;default:1" json:"payment_status"`                 // 支付状态(1:待支付 2:支付中 3:支付成功 4:支付失败 5:已退款)
	PaymentType          string  `gorm:"column:payment_type;size:20;index;not null" json:"payment_type"`                       // 支付类型(个人/企业)
	PaymentOrderNo       string  `gorm:"column:payment_order_no;size:64;uniqueIndex;not null" json:"payment_order_no"`         // 支付订单号，平台生成的唯一编号
	PaymentTransactionID string  `gorm:"column:payment_transaction_id;size:64;index" json:"payment_transaction_id"`            // 第三方交易ID
	PaymentTitle         string  `gorm:"column:payment_title;size:255;not null" json:"payment_title"`                          // 支付标题
	PaymentDescription   string  `gorm:"column:payment_description;size:1000" json:"payment_description"`                      // 支付描述
	NotifyURL            string  `gorm:"column:notify_url;size:255" json:"notify_url"`                                         // 支付回调通知地址

	PaymenterName   string          `gorm:"column:paymenter_name;size:100" json:"paymenter_name"`      // 付款人姓名/企业名称
	PaymenterEmail  string          `gorm:"column:paymenter_email;size:100" json:"paymenter_email"`    // 付款人邮箱
	PaymenterPhone  string          `gorm:"column:paymenter_phone;size:20" json:"paymenter_phone"`     // 付款人电话
	PaymentInfo     common.JSON     `gorm:"column:payment_info;type:json" json:"payment_info"`         // 账单信息(发票信息等)
	CompanyInfo     common.JSON     `gorm:"column:company_info;type:json" json:"company_info"`         // 企业付款信息(营业执照/税号等)
	CallbackData    common.JSON     `gorm:"column:callback_data;type:json" json:"callback_data"`       // 支付回调数据
	PlatformOptions common.JSON     `gorm:"column:platform_options;type:json" json:"platform_options"` // 支付平台特定配置
	PaymentTime     utils.MySQLTime `gorm:"column:payment_time" json:"payment_time"`                   // 支付成功时间
	ExpireTime      utils.MySQLTime `gorm:"column:expire_time;" json:"expire_time"`                    // 支付过期时间
	RefundInfo      common.JSON     `gorm:"column:refund_info;type:json" json:"refund_info"`           // 退款信息

	CreatedAt utils.MySQLTime `gorm:"column:created_at;not null" json:"created_at"` // 记录创建时间
	UpdatedAt utils.MySQLTime `gorm:"column:updated_at;not null" json:"updated_at"` // 记录更新时间
	DeletedAt gorm.DeletedAt  `gorm:"column:deleted_at" json:"deleted_at"`          // 软删除时间
}

// TableName 表名
func (Payment) TableName() string {
	return "payments"
}

// BeforeCreate 在创建记录前自动设置时间
func (payment *Payment) BeforeCreate(tx *gorm.DB) error {
	payment.CreatedAt = utils.MySQLTime(utils.GetTime())
	payment.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}

// BeforeUpdate 在更新记录前自动设置更新时间
func (payment *Payment) BeforeUpdate(tx *gorm.DB) error {
	payment.UpdatedAt = utils.MySQLTime(utils.GetTime())
	return nil
}
