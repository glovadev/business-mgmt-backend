package models

type DashboardReportResponse struct {
	Data struct {
		GrossSales             string `json:"gross_sales"`
		NetSales               string `json:"net_sales"`
		CashCollected          string `json:"cash_collected"`
		OutstandingReceivables string `json:"outstanding_receivables"`
		OverdueAmount          string `json:"overdue_amount"`
		COGS                   string `json:"cogs"`
		GrossProfit            string `json:"gross_profit"`
		StoreExpenses          string `json:"store_expenses"`
		CompanyExpenses        string `json:"company_expenses"`
		NetProfit              string `json:"net_profit"`
		InventoryValue         string `json:"inventory_value"`
		YellowCustomers        int    `json:"yellow_customers"`
		RedCustomers           int    `json:"red_customers"`
	} `json:"data"`
}

type GenericReportResponse struct {
	Data []map[string]interface{} `json:"data"`
}

type ProfitReportResponse struct {
	Data struct {
		Revenue     string `json:"revenue"`
		COGS        string `json:"cogs"`
		GrossProfit string `json:"gross_profit"`
		Expenses    string `json:"expenses"`
		NetProfit   string `json:"net_profit"`
	} `json:"data"`
}
