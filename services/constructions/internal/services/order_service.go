package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/pachv/constructions/constructions/internal/repositories"
)

type OrderService struct {
	repo   *repositories.OrderRepository
	mail   *MailSendingService
	logger *slog.Logger

	To []string // куда отправлять заказ (менеджерам)
}

func NewOrderService(repo *repositories.OrderRepository, mail *MailSendingService, logger *slog.Logger, to []string) *OrderService {
	return &OrderService{repo: repo, mail: mail, logger: logger, To: to}
}

type CartItem struct {
	ProductID string `json:"productId"`
	Qty       int    `json:"qty"`
}

type CreateOrderDTO struct {
	CompanyName     *string    `json:"companyName"`
	Email           string     `json:"email"`
	DeliveryAddress *string    `json:"deliveryAddress"`
	Comment         *string    `json:"comment"`
	Name            string     `json:"name"`
	Phone           string     `json:"phone"`
	Consent         bool       `json:"consent"`
	Items           []CartItem `json:"items"`
}

type OrderEmailData struct {
	OrderID         string
	CreatedAt       string
	CompanyName     *string
	Email           string
	DeliveryAddress *string
	Comment         *string
	Name            string
	Phone           string
	Items           []repositories.OrderItemForEmail
	Total           int
}

func (s *OrderService) CreateOrder(ctx context.Context, userID *string, dto CreateOrderDTO, templatePath string) (string, error) {
	if !dto.Consent {
		return "", fmt.Errorf("consent is required")
	}
	if dto.Email == "" || dto.Name == "" || dto.Phone == "" {
		return "", fmt.Errorf("email/name/phone are required")
	}
	if len(dto.Items) == 0 {
		return "", fmt.Errorf("items are required")
	}

	orderID := uuid.NewString()
	now := time.Now()

	items := make([]repositories.OrderItemInput, 0, len(dto.Items))
	for _, it := range dto.Items {
		items = append(items, repositories.OrderItemInput{
			ProductID: it.ProductID,
			Qty:       it.Qty,
		})
	}

	created, err := s.repo.CreateOrder(ctx, repositories.CreateOrderInput{
		OrderID:         orderID,
		UserID:          userID,
		CompanyName:     dto.CompanyName,
		Email:           dto.Email,
		DeliveryAddress: dto.DeliveryAddress,
		Comment:         dto.Comment,
		CustomerName:    dto.Name,
		CustomerPhone:   dto.Phone,
		Consent:         dto.Consent,
		CreatedAt:       now,
		Items:           items,
	})
	if err != nil {
		s.logger.Error("OrderService: CreateOrder failed", "err", err, "order_id", orderID)
		return "", err
	}

	subject := fmt.Sprintf("Новый заказ #%s", created.OrderID)

	data := OrderEmailData{
		OrderID:         created.OrderID,
		CreatedAt:       now.Format("2006-01-02 15:04:05"),
		CompanyName:     dto.CompanyName,
		Email:           dto.Email,
		DeliveryAddress: dto.DeliveryAddress,
		Comment:         dto.Comment,
		Name:            dto.Name,
		Phone:           dto.Phone,
		Items:           created.Items,
		Total:           created.Total,
	}

	if err := s.mail.SendHTMLFromTemplate(s.To, subject, templatePath, data); err != nil {
		s.logger.Error("OrderService: send email failed", "err", err, "order_id", created.OrderID)
		return "", fmt.Errorf("send email: %w", err)
	}

	return created.OrderID, nil
}
