package controllers

import (
	"API-Books/initializer"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ─── Infaq summaries ─────────────────────────────────────────────────────────

// GetInfaqSummaryByWeek godoc
// @Summary Infaq totals grouped by RT and week
// @Tags summary
// @Produce json
// @Param year query int false "Year (defaults to current year)"
// @Success 200 {array} map[string]interface{}
// @Router /api/summary/infaq/week [get]
func GetInfaqSummaryByWeek(ctx *gin.Context) {
	year := currentOrQueryYear(ctx)

	type Result struct {
		NeighborhoodID string  `json:"neighborhood_id"`
		Week           int     `json:"week"`
		Year           int     `json:"year"`
		TotalInfaq     float64 `json:"total_infaq"`
	}

	var results []Result
	err := initializer.DB.Table("infaqs").
		Select(`
            neighborhood_id,
            EXTRACT(WEEK FROM donated_at)::int AS week,
            EXTRACT(YEAR FROM donated_at)::int AS year,
            SUM(amount) AS total_infaq`).
		Where("EXTRACT(YEAR FROM donated_at) = ? AND deleted_at IS NULL", year).
		Group("neighborhood_id, week, year").
		Order("year, week").
		Scan(&results).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch summary"})
		return
	}
	ctx.JSON(http.StatusOK, results)
}

// GetInfaqSummaryByMonth godoc
// @Summary Infaq totals grouped by RT and month
// @Tags summary
// @Produce json
// @Param year query int false "Year (defaults to current year)"
// @Success 200 {array} map[string]interface{}
// @Router /api/summary/infaq/month [get]
func GetInfaqSummaryByMonth(ctx *gin.Context) {
	year := currentOrQueryYear(ctx)

	type Result struct {
		NeighborhoodID string  `json:"neighborhood_id"`
		Month          int     `json:"month"`
		Year           int     `json:"year"`
		TotalInfaq     float64 `json:"total_infaq"`
	}

	var results []Result
	err := initializer.DB.Table("infaqs").
		Select(`
            neighborhood_id,
            EXTRACT(MONTH FROM donated_at)::int AS month,
            EXTRACT(YEAR FROM donated_at)::int AS year,
            SUM(amount) AS total_infaq`).
		Where("EXTRACT(YEAR FROM donated_at) = ? AND deleted_at IS NULL", year).
		Group("neighborhood_id, month, year").
		Order("year, month").
		Scan(&results).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch summary"})
		return
	}
	ctx.JSON(http.StatusOK, results)
}

// ─── Loan summary ─────────────────────────────────────────────────────────────

// GetLoanMonthlySummary godoc
// @Summary Loan activity grouped by month for a given year
// @Tags summary
// @Produce json
// @Param year query int false "Year (defaults to current year)"
// @Success 200 {object} map[string]interface{}
// @Router /api/summary/loans/month [get]
func GetLoanMonthlySummary(ctx *gin.Context) {
	year := currentOrQueryYear(ctx)

	// New loans per month
	type MonthRow struct {
		Month      int     `json:"month"`
		NewLoans   float64 `json:"new_loans"`
		Repayments float64 `json:"repayments"`
	}

	var loanRows []MonthRow
	initializer.DB.Table("loans").
		Select(`EXTRACT(MONTH FROM start_date)::int AS month, SUM(amount) AS new_loans, 0 AS repayments`).
		Where("EXTRACT(YEAR FROM start_date) = ? AND deleted_at IS NULL", year).
		Group("month").
		Order("month").
		Scan(&loanRows)

	var payRows []MonthRow
	initializer.DB.Table("loan_payments").
		Select(`EXTRACT(MONTH FROM paid_at)::int AS month, 0 AS new_loans, SUM(amount) AS repayments`).
		Where("EXTRACT(YEAR FROM paid_at) = ? AND deleted_at IS NULL", year).
		Group("month").
		Order("month").
		Scan(&payRows)

	// Merge into 12-month array
	months := make([]MonthRow, 12)
	for i := range months {
		months[i].Month = i + 1
	}
	for _, r := range loanRows {
		if r.Month >= 1 && r.Month <= 12 {
			months[r.Month-1].NewLoans = r.NewLoans
		}
	}
	for _, r := range payRows {
		if r.Month >= 1 && r.Month <= 12 {
			months[r.Month-1].Repayments = r.Repayments
		}
	}

	// Totals
	type Summary struct {
		TotalLoans      float64    `json:"total_loans"`
		TotalRepayments float64    `json:"total_repayments"`
		TotalBorrowers  int64      `json:"total_borrowers"`
		Outstanding     float64    `json:"outstanding"`
		Months          []MonthRow `json:"months"`
	}

	var totalLoans, totalRepayments, outstanding float64
	var totalBorrowers int64
	initializer.DB.Table("loans").
		Select("COALESCE(SUM(amount),0), COALESCE(SUM(total_amount_paid),0), COALESCE(SUM(rest_payment),0)").
		Where("EXTRACT(YEAR FROM start_date) = ? AND deleted_at IS NULL", year).
		Row().Scan(&totalLoans, &totalRepayments, &outstanding)
	initializer.DB.Table("loans").
		Where("EXTRACT(YEAR FROM start_date) = ? AND deleted_at IS NULL", year).
		Count(&totalBorrowers)

	ctx.JSON(http.StatusOK, Summary{
		TotalLoans:      totalLoans,
		TotalRepayments: totalRepayments,
		TotalBorrowers:  totalBorrowers,
		Outstanding:     outstanding,
		Months:          months,
	})
}

// ─── Fund allocation ──────────────────────────────────────────────────────────

// GetFundAllocation godoc
// @Summary Fund allocation breakdown by purpose (DonationCuts)
// @Tags summary
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /api/summary/fund-allocation [get]
func GetFundAllocation(ctx *gin.Context) {
	type Result struct {
		Purpose string  `json:"purpose"`
		Total   float64 `json:"total"`
	}

	var results []Result
	err := initializer.DB.Table("donation_cuts").
		Select("purpose, SUM(amount) AS total").
		Where("deleted_at IS NULL").
		Group("purpose").
		Scan(&results).Error

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch fund allocation"})
		return
	}
	ctx.JSON(http.StatusOK, results)
}

// ─── helper ───────────────────────────────────────────────────────────────────

func currentOrQueryYear(ctx *gin.Context) int {
	if y := ctx.Query("year"); y != "" {
		if n, err := strconv.Atoi(y); err == nil {
			return n
		}
	}
	return time.Now().Year()
}
