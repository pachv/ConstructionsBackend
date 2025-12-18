package entity

import "time"

type OrderItem struct {
	ID                 string
	ProductID          string
	Qty                int
	ProductTitle       string
	ProductPrice       int
	ProductOldPrice    int
	ProductSalePercent int
	ProductImagePath   *string
}

type Order struct {
	ID              string
	UserID          *string
	CompanyName     *string
	Email           string
	DeliveryAddress *string
	Comment         *string
	CustomerName    string
	CustomerPhone   string
	Consent         bool
	CreatedAt       *time.Time
	Items           []OrderItem
}
