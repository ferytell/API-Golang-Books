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
    Loans          []Loan `gorm:"foreignKey:VillagerID"`
    Infaqs         []Infaq   `gorm:"foreignKey:VillagerID"`

}

type Infaq struct {
    gorm.Model
    VillagerID uint    `json:"villager_id"` // who donated
    Amount     float64 `json:"amount" validate:"required"`
    CollectedAt time.Time `json:"collected_at"`
    Cuts         []DonationCut `gorm:"foreignKey:InfaqID"`

}

type DonationCut struct {
    gorm.Model
    InfaqID uint
    Purpose string  // e.g. "transportation", "admission"
    Amount  float64
}


type Loan struct {
    gorm.Model
    VillagerID uint    `json:"villager_id"` // borrower
    Amount     float64 `json:"amount" validate:"required"`
    StartDate  time.Time `json:"start_date"`
    ActualEndDate *time.Time `json:"actual_end_date"`
    PlannedEndDate   time.Time `json:"planned_end_date"`
    TotalAmountPaid float64   `json:"total_amount_paid"`
    CurrentAmountPaid float64   `json:"current_amount_paid"`
    RestPayment    float64   `json:"rest_payment"`
    Reason         string    `json:"reason"` // optional
    Status     string    `json:"status"` // "ongoing", "paid", "defaulted"
}

type LoanPayment struct {
    gorm.Model
    LoanID   uint      `json:"loan_id"`
    Amount   float64   `json:"amount"`
    PaidAt   time.Time `json:"paid_at"`
}

type Committee struct {
    gorm.Model
    Name     string
    Position string // e.g. "Treasurer", "Secretary"
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
