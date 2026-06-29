package models

type DashboardStats struct {
	TotalProducts   int64 `json:"total_products"`
	TodayOrders     int64 `json:"today_orders"`
	ActiveStores    int64 `json:"active_stores"`
	StoresThreshold int64 `json:"stores_threshold"`
}
