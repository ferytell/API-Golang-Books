package models

import (
	"time"

	"gorm.io/gorm"
)

type Villager struct {
    gorm.Model
    Name           string `json:"name" validate:"required"`
    FamilyHeadName string `json:"family_head_name"`
    NeighborhoodID string `json:"neighborhood_id"` // e.g. "01", "02", "03"
}

type Infaq struct {
    gorm.Model
    VillagerID uint    `json:"villager_id"` // who donated
    Amount     float64 `json:"amount" validate:"required"`
    CollectedAt time.Time `json:"collected_at"`
}

type Loan struct {
    gorm.Model
    VillagerID uint    `json:"villager_id"` // borrower
    Amount     float64 `json:"amount" validate:"required"`
    StartDate  time.Time `json:"start_date"`
    EndDate    time.Time `json:"end_date"`
    Status     string    `json:"status"` // "ongoing", "paid", "defaulted"
}


type Repayment struct {
    gorm.Model
    LoanID   uint    `json:"loan_id"`
    Amount   float64 `json:"amount" validate:"required"`
    PaidAt   time.Time `json:"paid_at"`
}

type Fund struct {
    gorm.Model
    TotalBalance float64 `json:"total_balance"`
    LastUpdated  time.Time `json:"last_updated"`
}
