package dto

import "time"

type BaseOrder struct {
	ID             string    `gorm:"primaryKey;size:100;not null" json:"id"`
	OrderID        string    `gorm:"size:100;not null" json:"order_id"`
	Priority       string    `gorm:"size:100;not null" json:"priority"`
	Status         string    `gorm:"size:100;not null" json:"status"`
	ProcessingTime int       `gorm:"not null;default:0" json:"processing_time"`
	Lock           bool      `gorm:"not null;default:0" json:"lock"`
	CreatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (BaseOrder) TableName() string {
	return "orders"
}
