package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pachv/constructions/constructions/internal/domain/entity"
)

type AdminDashboardService struct {
	db *sqlx.DB
}

func NewAdminDashboardService(db *sqlx.DB) *AdminDashboardService {
	return &AdminDashboardService{db: db}
}

func (s *AdminDashboardService) GetStats(ctx context.Context) (entity.AdminDashboardStats, error) {
	const q = `
		SELECT
			(SELECT COALESCE(COUNT(*), 0) FROM users)  AS total_users,
			(SELECT COALESCE(COUNT(*), 0) FROM orders) AS total_orders,
			(SELECT COALESCE(COUNT(*), 0)
			 FROM orders
			 WHERE created_at IS NOT NULL
			   AND created_at >= date_trunc('day', now())
			) AS orders_today,
			(SELECT COALESCE(COUNT(*), 0)
			 FROM orders
			 WHERE created_at IS NOT NULL
			   AND created_at >= (now() - interval '1 month')
			) AS orders_last_month
	`

	var out entity.AdminDashboardStats
	if err := s.db.GetContext(ctx, &out, q); err != nil {
		return out, fmt.Errorf("admin dashboard stats: %w", err)
	}
	return out, nil
}
