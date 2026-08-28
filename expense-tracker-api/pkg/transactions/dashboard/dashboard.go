package dashboard

import (
	"net/http"
	"expense-tracker-api/pkg/response"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardSummary struct {
	TotalIncome      float64 `json:"total_income"`
	TotalExpense     float64 `json:"total_expense"`
	Balance          float64 `json:"balance"`
	ThisMonthIncome  float64 `json:"this_month_income"`
	ThisMonthExpense float64 `json:"this_month_expense"`
	LastMonthIncome  float64 `json:"last_month_income"`
	LastMonthExpense float64 `json:"last_month_expense"`
	Period           string  `json:"period"`
}

type ExpenseBreakdown struct {
	ByCategory map[string]float64 `json:"by_category,omitempty"`
	ByMember   map[string]float64 `json:"by_member,omitempty"`
	ByPayment  map[string]float64 `json:"by_payment_method,omitempty"`
}

type DashboardData struct {
	Summary    DashboardSummary   `json:"summary"`
	Breakdown  ExpenseBreakdown   `json:"breakdown"`
}

type DashboardHandler struct {
	DB *pgxpool.Pool
}

func NewDashboardHandler(db *pgxpool.Pool) *DashboardHandler {
	return &DashboardHandler{DB: db}
}

func (h *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "monthly"
	}

	data := DashboardData{
		Summary: DashboardSummary{
			TotalIncome:      0,
			TotalExpense:     0,
			Balance:          0,
			ThisMonthIncome:  0,
			ThisMonthExpense: 0,
			LastMonthIncome:  0,
			LastMonthExpense: 0,
			Period:           period,
		},
		Breakdown: ExpenseBreakdown{
			ByCategory:  make(map[string]float64),
			ByMember:    make(map[string]float64),
			ByPayment:   make(map[string]float64),
		},
	}

	response.JSON(w, http.StatusOK, response.HTTPResponse{
		Success: true,
		Data:    data,
	})
}