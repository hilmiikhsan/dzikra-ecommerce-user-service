package dto

type GetDashboardResponse struct {
	TotalAmount         int `json:"total_amount"`
	TotalExpenses       int `json:"total_expenses"`
	TotalTransaction    int `json:"total_transaction"`
	TotalSellingProduct int `json:"total_selling_product"`
	TotalCapital        int `json:"total_capital"`
	Netsales            int `json:"net_sales"`
	ProfitLoss          int `json:"profit_loss"`
}
