package controllers

import (
	"API-Books/initializer"
	"API-Books/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// swagger:model LoanRequest
type loanRequest struct {
	VillagerID     uint    `json:"villager_id"      binding:"required"`
	Amount         float64 `json:"amount"           binding:"required,gt=0"`
	PlannedEndDate string  `json:"planned_end_date" binding:"required"` // "YYYY-MM-DD"
	Reason         string  `json:"reason"`
	Notes          string  `json:"notes"`
}
type paymentRequest struct {
	Amount    float64 `json:"amount"     binding:"required,gt=0"`
	PaidAt    string  `json:"paid_at"`    // "YYYY-MM-DD", defaults to today
	WeekLabel string  `json:"week_label"` // e.g. "Minggu 1"
	Notes     string  `json:"notes"`
}

// ─── Loans ────────────────────────────────────────────────────────────────────

// RequestLoan godoc
// @Summary Create a new loan
// @Tags loans
// @Accept json
// @Produce json
// @Param loan body loanRequest true "Loan data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/loans [post]
func RequestLoan(ctx *gin.Context) {
	var req loanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var villager models.Villager
	if err := initializer.DB.First(&villager, req.VillagerID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "villager not found"})
		return
	}

	plannedEndDate, err := time.Parse("2006-01-02", req.PlannedEndDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	loan := models.Loan{
		VillagerID:        req.VillagerID,
		Amount:            req.Amount,
		StartDate:         time.Now(),
		PlannedEndDate:    plannedEndDate,
		TotalAmountPaid:   0,
		CurrentAmountPaid: 0,
		RestPayment:       req.Amount,
		Reason:            req.Reason,
		Notes:             req.Notes,
		Status:            "ongoing",
	}

	if err := initializer.DB.Create(&loan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create loan"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "loan created successfully", "loan": loan})
}

// GetAllLoans godoc
// @Summary Get all loans
// @Tags loans
// @Produce json
// @Param status       query string false "Filter by status: ongoing|paid|defaulted"
// @Param villager_id  query int    false "Filter by villager"
// @Param with_villager query bool  false "Preload villager info"
// @Success 200 {array} models.Loan
// @Failure 500 {object} map[string]string
// @Router /api/loans [get]
func GetAllLoans(ctx *gin.Context) {
	var loans []models.Loan
	q := initializer.DB.Model(&models.Loan{})

	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if vid := ctx.Query("villager_id"); vid != "" {
		q = q.Where("villager_id = ?", vid)
	}
	if ctx.Query("with_villager") == "true" {
		q = q.Preload("Villager")
	}

	if err := q.Order("created_at desc").Find(&loans).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch loans"})
		return
	}
	ctx.JSON(http.StatusOK, loans)
}

// GetLoanByID godoc
// @Summary Get a single loan by ID (includes payments)
// @Tags loans
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {object} models.Loan
// @Failure 404 {object} map[string]string
// @Router /api/loans/{id} [get]
func GetLoanByID(ctx *gin.Context) {
	id := ctx.Param("id")
	var loan models.Loan
	if err := initializer.DB.Preload("Payments").Preload("Villager").First(&loan, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}
	ctx.JSON(http.StatusOK, loan)
}

// DeleteLoan godoc
// @Summary Delete a loan
// @Tags loans
// @Param id path int true "Loan ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/loans/{id} [delete]
func DeleteLoan(ctx *gin.Context) {
	id := ctx.Param("id")
	var loan models.Loan
	if err := initializer.DB.First(&loan, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}
	// also remove associated payments
	initializer.DB.Where("loan_id = ?", loan.ID).Delete(&models.LoanPayment{})
	initializer.DB.Delete(&loan)
	ctx.JSON(http.StatusOK, gin.H{"message": "loan deleted"})
}

// ─── Payments ─────────────────────────────────────────────────────────────────

// CreatePayment godoc
// @Summary Record a payment for a loan
// @Tags payments
// @Accept json
// @Produce json
// @Param id   path int            true "Loan ID"
// @Param body body paymentRequest true "Payment data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/loans/{id}/payments [post]
func CreatePayment(c *gin.Context) {
	loanID := c.Param("id")

	var loan models.Loan
	if err := initializer.DB.First(&loan, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	if loan.Status != "ongoing" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "loan is not ongoing"})
		return
	}

	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount > loan.RestPayment+0.01 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment exceeds remaining balance"})
		return
	}

	paidAt := time.Now()
	if req.PaidAt != "" {
		if t, err := time.Parse("2006-01-02", req.PaidAt); err == nil {
			paidAt = t
		}
	}

	payment := models.LoanPayment{
		LoanID:    loan.ID,
		Amount:    req.Amount,
		PaidAt:    paidAt,
		WeekLabel: req.WeekLabel,
		Notes:     req.Notes,
	}

	if err := initializer.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record payment"})
		return
	}

	// Update loan aggregates
	loan.TotalAmountPaid += req.Amount
	loan.CurrentAmountPaid += req.Amount
	loan.RestPayment = loan.Amount - loan.TotalAmountPaid
	if loan.RestPayment <= 0.001 {
		loan.RestPayment = 0
		loan.Status = "paid"
		now := time.Now()
		loan.ActualEndDate = &now
	}
	initializer.DB.Save(&loan)

	c.JSON(http.StatusCreated, gin.H{
		"message": "payment recorded",
		"payment": payment,
		"loan":    loan,
	})
}

// GetPaymentsByLoan godoc
// @Summary Get all payments for a loan
// @Tags payments
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {array} models.LoanPayment
// @Failure 404 {object} map[string]string
// @Router /api/loans/{id}/payments [get]
func GetPaymentsByLoan(c *gin.Context) {
	loanID := c.Param("id")
	var payments []models.LoanPayment
	if err := initializer.DB.Where("loan_id = ?", loanID).Order("paid_at asc").Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// GetAllPayments godoc
// @Summary Get all payments across all loans
// @Tags payments
// @Produce json
// @Success 200 {array} models.LoanPayment
// @Router /api/loan_payments [get]
func GetAllPayments(c *gin.Context) {
	var payments []models.LoanPayment
	if err := initializer.DB.Order("paid_at desc").Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// UpdatePayment godoc
// @Summary Update a payment record
// @Tags payments
// @Accept json
// @Produce json
// @Param id   path int            true "Payment ID"
// @Param body body paymentRequest true "Updated payment"
// @Success 200 {object} models.LoanPayment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/loan_payments/{id} [put]
func UpdatePayment(c *gin.Context) {
	id := c.Param("id")
	var payment models.LoanPayment
	if err := initializer.DB.First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	var req paymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Roll back old amount on loan
	var loan models.Loan
	if err := initializer.DB.First(&loan, payment.LoanID).Error; err == nil {
		loan.TotalAmountPaid -= payment.Amount
		loan.CurrentAmountPaid -= payment.Amount

		// Apply new amount
		loan.TotalAmountPaid += req.Amount
		loan.CurrentAmountPaid += req.Amount
		loan.RestPayment = loan.Amount - loan.TotalAmountPaid
		if loan.RestPayment <= 0.001 {
			loan.RestPayment = 0
			loan.Status = "paid"
		} else {
			loan.Status = "ongoing"
			loan.ActualEndDate = nil
		}
		initializer.DB.Save(&loan)
	}

	payment.Amount = req.Amount
	payment.WeekLabel = req.WeekLabel
	payment.Notes = req.Notes
	if req.PaidAt != "" {
		if t, err := time.Parse("2006-01-02", req.PaidAt); err == nil {
			payment.PaidAt = t
		}
	}
	initializer.DB.Save(&payment)

	c.JSON(http.StatusOK, payment)
}

// DeletePayment godoc
// @Summary Delete a payment and revert loan balance
// @Tags payments
// @Param id path int true "Payment ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/loan_payments/{id} [delete]
func DeletePayment(c *gin.Context) {
	id := c.Param("id")
	var payment models.LoanPayment
	if err := initializer.DB.First(&payment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
		return
	}

	// Revert loan aggregates
	var loan models.Loan
	if err := initializer.DB.First(&loan, payment.LoanID).Error; err == nil {
		loan.TotalAmountPaid -= payment.Amount
		loan.CurrentAmountPaid -= payment.Amount
		loan.RestPayment = loan.Amount - loan.TotalAmountPaid
		if loan.RestPayment > 0 {
			loan.Status = "ongoing"
			loan.ActualEndDate = nil
		}
		initializer.DB.Save(&loan)
	}

	initializer.DB.Delete(&payment)
	c.JSON(http.StatusOK, gin.H{"message": "payment deleted"})
}

// PaymentLoan is kept for backward compat with old batch route
// Prefer POST /api/loans/:id/payments instead
func PaymentLoan(ctx *gin.Context) {
	ctx.JSON(http.StatusGone, gin.H{"error": "deprecated — use POST /api/loans/:id/payments"})
}
