package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type OrderRepository struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewOrderRepository(db *sqlx.DB, logger *slog.Logger) *OrderRepository {
	return &OrderRepository{db: db, logger: logger}
}

// минимальный срез полей товара из БД
type catalogProductRow struct {
	ID          string        `db:"id"`
	Title       string        `db:"title"`
	Price       int           `db:"price"`
	OldPrice    sql.NullInt64 `db:"old_price"`
	SalePercent sql.NullInt64 `db:"sale_percent"`
	ImagePath   *string       `db:"image_path"`
}

type OrderItemInput struct {
	ProductID string
	Qty       int
}

type CreateOrderInput struct {
	OrderID         string
	UserID          *string
	CompanyName     *string
	Email           string
	DeliveryAddress *string
	Comment         *string
	CustomerName    string
	CustomerPhone   string
	Consent         bool
	CreatedAt       time.Time

	Items []OrderItemInput
}

type OrderItemForEmail struct {
	ProductID   string
	Title       string
	Qty         int
	Price       int
	OldPrice    int
	SalePercent int
	ImagePath   *string
	LineTotal   int
}

type CreatedOrder struct {
	OrderID string
	Items   []OrderItemForEmail
	Total   int
}

func (r *OrderRepository) CreateOrder(ctx context.Context, in CreateOrderInput) (CreatedOrder, error) {
	r.logger.Debug("CreateOrder: start", "order_id", in.OrderID)

	if len(in.Items) == 0 {
		return CreatedOrder{}, fmt.Errorf("CreateOrder: empty items")
	}
	for _, it := range in.Items {
		if strings.TrimSpace(it.ProductID) == "" || it.Qty <= 0 {
			return CreatedOrder{}, fmt.Errorf("CreateOrder: invalid item productId/qty")
		}
	}

	// 1) получить товары по id
	ids := make([]string, 0, len(in.Items))
	qtyByID := make(map[string]int, len(in.Items))
	for _, it := range in.Items {
		ids = append(ids, it.ProductID)
		qtyByID[it.ProductID] += it.Qty
	}

	const productsQ = `
	SELECT id, title, price, old_price, sale_percent, image_path
	FROM catalog_products
	WHERE id = ANY($1::text[])
`

	var products []catalogProductRow
	if err := r.db.SelectContext(ctx, &products, productsQ, pq.Array(ids)); err != nil {
		return CreatedOrder{}, fmt.Errorf("select products: %w", err)
	}
	if len(products) == 0 {
		return CreatedOrder{}, fmt.Errorf("select products: no products found")
	}

	// проверка: все ли id найдены
	found := map[string]bool{}
	for _, p := range products {
		found[p.ID] = true
	}
	for _, id := range qtyByID {
		_ = id
	}
	for pid := range qtyByID {
		if !found[pid] {
			return CreatedOrder{}, fmt.Errorf("product not found: %s", pid)
		}
	}

	// 2) транзакция: orders + order_items
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return CreatedOrder{}, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	const insertOrderQ = `
		INSERT INTO orders (
			id, user_id,
			company_name, email, delivery_address, comment,
			customer_name, customer_phone, consent,
			created_at
		) VALUES (
			:id, :user_id,
			:company_name, :email, :delivery_address, :comment,
			:customer_name, :customer_phone, :consent,
			:created_at
		)
	`

	orderArgs := map[string]any{
		"id":               in.OrderID,
		"user_id":          in.UserID,
		"company_name":     in.CompanyName,
		"email":            in.Email,
		"delivery_address": in.DeliveryAddress,
		"comment":          in.Comment,
		"customer_name":    in.CustomerName,
		"customer_phone":   in.CustomerPhone,
		"consent":          in.Consent,
		"created_at":       in.CreatedAt,
	}

	if _, err := tx.NamedExecContext(ctx, insertOrderQ, orderArgs); err != nil {
		r.logger.Error("CreateOrder: insert order failed", "err", err, "query", insertOrderQ, "order_id", in.OrderID)
		return CreatedOrder{}, fmt.Errorf("insert order: %w", err)
	}

	const insertItemQ = `
		INSERT INTO order_items (
			id, order_id,
			product_id, qty,
			product_title, product_price, product_old_price, product_sale_percent, product_image_path,
			created_at
		) VALUES (
			:id, :order_id,
			:product_id, :qty,
			:product_title, :product_price, :product_old_price, :product_sale_percent, :product_image_path,
			:created_at
		)
	`

	emailItems := make([]OrderItemForEmail, 0, len(products))
	total := 0

	for _, p := range products {
		qty := qtyByID[p.ID]

		oldPrice := p.Price
		if p.OldPrice.Valid {
			oldPrice = int(p.OldPrice.Int64)
		}
		salePercent := 0
		if p.SalePercent.Valid {
			salePercent = int(p.SalePercent.Int64)
		}

		line := qty * p.Price
		total += line

		itemID := fmt.Sprintf("oi-%s-%d", in.OrderID, len(emailItems)+1) // простой id; лучше uuid в реале

		itemArgs := map[string]any{
			"id":                   itemID,
			"order_id":             in.OrderID,
			"product_id":           p.ID,
			"qty":                  qty,
			"product_title":        p.Title,
			"product_price":        p.Price,
			"product_old_price":    oldPrice,
			"product_sale_percent": salePercent,
			"product_image_path":   p.ImagePath,
			"created_at":           in.CreatedAt,
		}

		if _, err := tx.NamedExecContext(ctx, insertItemQ, itemArgs); err != nil {
			r.logger.Error("CreateOrder: insert item failed", "err", err, "query", insertItemQ, "order_id", in.OrderID, "product_id", p.ID)
			return CreatedOrder{}, fmt.Errorf("insert order item: %w", err)
		}

		emailItems = append(emailItems, OrderItemForEmail{
			ProductID:   p.ID,
			Title:       p.Title,
			Qty:         qty,
			Price:       p.Price,
			OldPrice:    oldPrice,
			SalePercent: salePercent,
			ImagePath:   p.ImagePath,
			LineTotal:   line,
		})
	}

	if err := tx.Commit(); err != nil {
		return CreatedOrder{}, fmt.Errorf("commit: %w", err)
	}

	r.logger.Debug("CreateOrder: done", "order_id", in.OrderID, "total", total, "items", len(emailItems))
	return CreatedOrder{OrderID: in.OrderID, Items: emailItems, Total: total}, nil
}
