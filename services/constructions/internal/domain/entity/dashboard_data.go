package entity

type AdminDashboardStats struct {
	TotalUsers      int `db:"total_users" json:"totalUsers"`
	TotalOrders     int `db:"total_orders" json:"totalOrders"`
	OrdersToday     int `db:"orders_today" json:"ordersToday"`
	OrdersLastMonth int `db:"orders_last_month" json:"ordersLastMonth"`
}
