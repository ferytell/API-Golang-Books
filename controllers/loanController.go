package controllers

import (
	"API-Books/initializer"
	"API-Books/models"
	"API-Books/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)





type LoanRequest struct {
    VillagerID     uint    `json:"villager_id" binding:"required"`
    Amount         float64 `json:"amount" binding:"required"`
    PlannedEndDate string  `json:"planned_end_date" binding:"required"` // ISO date
    Reason         string  `json:"reason"`
}

// RequestLoan godoc
// @Summary Request a new loan
// @Description Villager requests a new loan
// @Tags loans
// @Accept json
// @Produce json
// @Param loan body LoanRequest true "Loan Request Data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /loans [post]
func RequestLoan(ctx *gin.Context) {
    var req LoanRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // check villager exists
    var villager models.Villager
    if err := initializer.DB.First(&villager, req.VillagerID).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "villager not found"})
        return
    }

    // parse planned end date
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
        Status:            "ongoing",
    }

    if err := initializer.DB.Create(&loan).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create loan"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{
        "message": "loan request created successfully",
        "loan":    loan,
    })
}

// PaymentLoan godoc
// @Summary Make a payment towards a loan
// @Description Villager makes a payment towards their loan
// @Tags loans
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payment body map[string]float64 true "Payment Data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /loans/{id}/payment [post]
func PaymentLoan(ctx *gin.Context) {
	loanID := ctx.Param("id")

	var body struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var loan models.Loan
	if err := initializer.DB.First(&loan, loanID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "loan not found"})
		return
	}

	if loan.Status != "ongoing" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "loan is not ongoing"})
		return
	}

	if body.Amount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payment amount must be positive"})
		return
	}

	if body.Amount > loan.RestPayment {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "payment exceeds remaining balance"})
		return
	}

	loan.CurrentAmountPaid += body.Amount
	loan.TotalAmountPaid += body.Amount
	loan.RestPayment -= body.Amount

	if loan.RestPayment <= 0.0001 { // tolerance for float
		loan.RestPayment = 0
		loan.Status = "paid"
		now := time.Now()
		loan.ActualEndDate = &now
	}

	if err := initializer.DB.Save(&loan).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update loan"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "payment successful",
		"loan":    loan,
	})
}

// GetVillagerLoans godoc
// @Summary Get loans for a specific villager
// @Description Retrieve all loans associated with a villager by their ID
// @Tags villagers
// @Produce json
// @Param id path int true "Villager ID"
// @Success 200 {array} models.Loan
// @Failure 404 {object} map[string]string
// @Router /villagers/{id}/loans [get]
func GetVillagerLoans(ctx *gin.Context) {
	id := ctx.Param("id")
	var villager models.Villager

	if err := initializer.DB.Preload("Loans").First(&villager, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Villager not found"})
		return
	}

	ctx.JSON(http.StatusOK, villager.Loans)
}

// GetAllLoans godoc
// @Summary Get all loans
// @Description Retrieve all loans in the system
// @Tags loans
// @Produce json
// @Success 200 {array} models.Loan
// @Failure 500 {object} map[string]string
// @Router /loans [get]
func GetAllLoans(ctx *gin.Context) {
	var loans []models.Loan
	if err := initializer.DB.Find(&loans).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch loans"})
		return
	}

	ctx.JSON(http.StatusOK, loans)
}




// CreatePayment godoc
// @Summary Create a loan payment
// @Description Record a new payment towards a loan
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Loan ID"
// @Param payment body models.LoanPayment true "Payment Data"
// @Success 201 {object} models.LoanPayment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /loan_payments/{id} [post]
func CreatePayment(c *gin.Context) {
    var payment models.LoanPayment
    loanID := c.Param("id")

    if err := c.ShouldBindJSON(&payment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    payment.LoanID = utils.StringToUint(loanID) // helper to convert string -> uint

    if payment.PaidAt.IsZero() {
        payment.PaidAt = time.Now()
    }

    if err := initializer.DB.Create(&payment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // update Loan aggregates
    var loan models.Loan
    if err := initializer.DB.First(&loan, payment.LoanID).Error; err == nil {
        loan.TotalAmountPaid += payment.Amount
        loan.CurrentAmountPaid += payment.Amount
        loan.RestPayment = loan.Amount - loan.TotalAmountPaid
        if loan.RestPayment <= 0 {
            loan.Status = "paid"
        }
        initializer.DB.Save(&loan)
    }

    c.JSON(http.StatusCreated, payment)
}

// UpdatePayment godoc
// @Summary Update a loan payment
// @Description Update details of a specific loan payment
// @Tags payments
// @Accept json
// @Produce json
// @Param id path int true "Payment ID"
// @Param payment body models.LoanPayment true "Payment Data"
// @Success 200 {object} models.LoanPayment
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /loan_payments/{id} [put]
func UpdatePayment(c *gin.Context) {
    id := c.Param("id")
    var payment models.LoanPayment
    if err := initializer.DB.First(&payment, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "payment not found"})
        return
    }

    var input models.LoanPayment
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // rollback loan totals first
    var loan models.Loan
    if err := initializer.DB.First(&loan, payment.LoanID).Error; err == nil {
        loan.TotalAmountPaid -= payment.Amount
        loan.CurrentAmountPaid -= payment.Amount
    }

    // update payment
    payment.Amount = input.Amount
    payment.PaidAt = input.PaidAt
    initializer.DB.Save(&payment)

    // reapply to loan
    loan.TotalAmountPaid += payment.Amount
    loan.CurrentAmountPaid += payment.Amount
    loan.RestPayment = loan.Amount - loan.TotalAmountPaid
    if loan.RestPayment <= 0 {
        loan.Status = "paid"
    } else {
        loan.Status = "ongoing"
    }
    initializer.DB.Save(&loan)

    c.JSON(http.StatusOK, payment)
}

// GetAllPayments godoc
// @Summary Get all loan payments
// @Description Retrieve all loan payments in the system
// @Tags payments
// @Produce json
// @Success 200 {array} models.LoanPayment
// @Failure 500 {object} map[string]string
// @Router /loan_payments [get]
func GetAllPayments(c *gin.Context) {
    var payments []models.LoanPayment

    if err := initializer.DB.Preload("Loan").Find(&payments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, payments)
}

// GET /loans/:id/payments

// GetPaymentsByLoan godoc
// @Summary Get payments for a specific loan
// @Description Retrieve all payments associated with a loan by its ID
// @Tags loans
// @Produce json
// @Param id path int true "Loan ID"
// @Success 200 {array} models.LoanPayment
// @Failure 404 {object} map[string]string
// @Router /loans/{id}/payments [get]
func GetPaymentsByLoan(c *gin.Context) {
    loanID := c.Param("id")
    var payments []models.LoanPayment

    if err := initializer.DB.Where("loan_id = ?", loanID).Find(&payments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, payments)
}


