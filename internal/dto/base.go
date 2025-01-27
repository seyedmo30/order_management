package dto

import "time"

type BaseOrder struct {
	ID             string     `gorm:"primaryKey;size:100;not null" json:"id"`
	OrderID        *string    `gorm:"size:100;not null;unique" json:"order_id"`
	Priority       *string    `gorm:"size:100;not null" json:"priority"`
	Status         *string    `gorm:"size:100;not null" json:"status"`
	ProcessingTime *int       `gorm:"not null;default:0" json:"processing_time"`
	Lock           *bool      `gorm:"not null;default:0" json:"lock"`
	CreatedAt      *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      *time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (BaseOrder) TableName() string {
	return "orders"
}



type BaseCreateOrderRequest struct {
	OrderID        string `json:"order_id" validate:"required"`
	Priority       string `json:"priority" validate:"required,oneof=High Normal"`
	ProcessingTime int    `json:"processing_time" validate:"required,min=1"`
}
