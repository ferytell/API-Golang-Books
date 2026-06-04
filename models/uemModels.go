package models

import (
	"time"

	"gorm.io/gorm"
)

// type Villager struct {
//     gorm.Model
//     Name           string `json:"name" validate:"required"`
//     FamilyHeadName string `json:"family_head_name"`
//     NeighborhoodID string `json:"neighborhood_id"` // e.g. "01", "02", "03"
//     Loans          []Loan `gorm:"foreignKey:VillagerID"`
//     Infaqs         []Infaq   `gorm:"foreignKey:VillagerID"`

// }

// ENHANCED: Villager with Infaq relationship
type Villager struct {
    gorm.Model
    Name           string    `json:"name" validate:"required"`
    FamilyHeadName string    `json:"family_head_name"`
    NeighborhoodID string    `json:"neighborhood_id"`
    InfaqNumber    string    `json:"infaq_number"` // e.g. "101023"
    IsActive       bool      `json:"is_active" gorm:"default:true"`
    Loans          []Loan    `gorm:"foreignKey:VillagerID"`
    Infaqs         []Infaq   `gorm:"foreignKey:VillagerID"`
}

// type Infaq struct {
//     gorm.Model
//     ID             uint      `json:"id" gorm:"primaryKey"`
//     VillagerID *uint     `json:"villager_id"` // who donated
//     NeighborhoodID string      `json:"neighborhood_id" gorm:"not null"`
//     Amount     float64 `json:"amount" validate:"required"`
//     DonatedAt      time.Time `json:"donated_at"`
//     CollectedAt time.Time `json:"collected_at"`
//     Cuts         []DonationCut `gorm:"foreignKey:InfaqID"`

// }

// ENHANCED: Infaq with weekly tracking
type Infaq struct {
    gorm.Model
    ID             uint      `json:"id" gorm:"primaryKey"`
    VillagerID     *uint     `json:"villager_id"`
    NeighborhoodID string    `json:"neighborhood_id" gorm:"not null"`
    Amount         float64   `json:"amount" validate:"required"`
    DonatedAt      time.Time `json:"donated_at"`
    WeekNumber     int       `json:"week_number"`      // 1-48
    Month          int       `json:"month"`            // 1-12
    Year           int       `json:"year"`
    CollectedAt    time.Time `json:"collected_at"`
    Cuts           []DonationCut `gorm:"foreignKey:InfaqID"`
	Notes          string    `json:"notes"`
}

// swagger:model Loan
type Loan struct {
	gorm.Model
	ID               uint      `json:"id" gorm:"primaryKey"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
    DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	VillagerID        uint       `json:"villager_id"          gorm:"not null;index"  validate:"required"`
	Amount            float64    `json:"amount"               gorm:"not null"        validate:"required,gt=0"`
	StartDate         time.Time  `json:"start_date"`
	PlannedEndDate    time.Time  `json:"planned_end_date"`
	ActualEndDate     *time.Time `json:"actual_end_date,omitempty"`
	TotalAmountPaid   float64    `json:"total_amount_paid"    gorm:"default:0"`
	CurrentAmountPaid float64    `json:"current_amount_paid"  gorm:"default:0"`
	RestPayment       float64    `json:"rest_payment"         gorm:"default:0"`
	Reason            string     `json:"reason"`
	Notes             string     `json:"notes"`
	// Status: "ongoing", "paid", "defaulted"
	Status            string     `json:"status"               gorm:"default:'ongoing';index"`

	// Associations (loaded only when Preloaded)
	Villager  *Villager    `json:"villager,omitempty"  gorm:"foreignKey:VillagerID"`
	Payments  []LoanPayment `json:"payments,omitempty" gorm:"foreignKey:LoanID"`
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
// NEW: Neighborhood (RT) model
type Neighborhood struct {
    gorm.Model
    Code        string      `json:"code" gorm:"unique"` // "01", "02", "03"
    Name        string      `json:"name"`
    Villagers   []Villager  `gorm:"foreignKey:NeighborhoodID"`
}



// NEW: CashBookEntry for monthly ledger tracking
type CashBookEntry struct {
    gorm.Model
    BookType    string    `json:"book_type"`    // "general", "dkm", "operational", "management"
    Month       int       `json:"month"`
    Year        int       `json:"year"`
    Date        time.Time `json:"date"`
    Description string    `json:"description"`
    Income      float64   `json:"income"`
    Expense     float64   `json:"expense"`
    Balance     float64   `json:"balance"`
}

// NEW: FundAllocation for automatic 10% splits
type FundAllocation struct {
    gorm.Model
    InfaqID     uint    `json:"infaq_id"`
    DKMPercent      float64 `json:"dkm_percent" gorm:"default:10"`
    OperationalPercent float64 `json:"operational_percent" gorm:"default:10"`
    ManagementPercent  float64 `json:"management_percent" gorm:"default:10"`
}

// NEW: LoanCard for individual borrower tracking
type LoanCard struct {
    gorm.Model
    LoanID      uint      `json:"loan_id"`
    VillagerID  uint      `json:"villager_id"`
    PrintDate   time.Time `json:"print_date"`
    GeneratedBy string    `json:"generated_by"`
}



// CashBookSummary represents one of the 4 cash books (General, DKM, Operational, Management)
type CashBookSummary struct {
	BookType    string  `json:"book_type"`    // "general", "dkm", "operational", "management"
	OpeningBalance float64 `json:"opening_balance"`
	Entries     []CashBookEntryDetail `json:"entries"`
	ClosingBalance float64 `json:"closing_balance"`
	TotalIncome float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
}

type CashBookEntryDetail struct {
	No          int     `json:"no"`
	Description string  `json:"description"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
	Balance     float64 `json:"balance"`
}

// MonthlyReport matches your MLAPORAN sheet structure
type MonthlyReport struct {
	Period          string          `json:"period"`          // "Januari 2025", "Februari 2025", etc.
	PeriodEnglish   string          `json:"period_english"`
	Year            int             `json:"year"`
	Month           int             `json:"month"`
	
	// I. Keadaan Kas Umum UEM (General Cash)
	GeneralCash     CashBookSummary `json:"general_cash"`
	
	// II. Keadaan Kas Infak Untuk DKM 10%
	DKMCash         CashBookSummary `json:"dkm_cash"`
	
	// III. Keadaan Kas Operasional UEM 10%
	OperationalCash CashBookSummary `json:"operational_cash"`
	
	// IV. Keadaan Kas Akomodasi Pengurus UEM 10%
	ManagementCash  CashBookSummary `json:"management_cash"`
	
	// Summary totals
	TotalAllBalance float64 `json:"total_all_balance"`
	
	// Loan statistics for the month
	LoanStats       MonthlyLoanStats `json:"loan_stats"`
	
	// Infaq statistics by RT
	InfaqByRT       []RTInfaqSummary `json:"infaq_by_rt"`
	
	GeneratedAt     time.Time `json:"generated_at"`
}

type MonthlyLoanStats struct {
	NewLoans        int64     `json:"new_loans"`
	TotalDisbursed  float64 `json:"total_disbursed"`
	TotalRepaid     float64 `json:"total_repaid"`
	ActiveLoans     int64     `json:"active_loans"`
	CompletedLoans  int64     `json:"completed_loans"`
}

type RTInfaqSummary struct {
	NeighborhoodID string  `json:"neighborhood_id"`
	TotalInfaq     float64 `json:"total_infaq"`
	ResidentCount  int64   `json:"resident_count"`
}

// SemesterReport for 6-month periods
type SemesterReport struct {
	Period          string  `json:"period"`          // "Januari - Juni 2025"
	Semester        int     `json:"semester"`        // 1 or 2
	Year            int     `json:"year"`
	
	// Buku Kas Umum summary
	TotalInfaqCollected   float64 `json:"total_infaq_collected"`
	TotalLoanRepayments   float64 `json:"total_loan_repayments"`
	TotalLoanDisbursed    float64 `json:"total_loan_disbursed"`
	TotalDKMTransferred   float64 `json:"total_dkm_transferred"`
	TotalOperationalUsed  float64 `json:"total_operational_used"`
	TotalManagementUsed   float64 `json:"total_management_used"`
	
	// Data Perputaran Keuangan (Financial Turnover)
	FinancialTurnover     SemesterFinancialTurnover `json:"financial_turnover"`
	
	// Data Pemanfaat UEM (Beneficiary Data)
	BeneficiaryStats      BeneficiaryStatistics `json:"beneficiary_stats"`
	
	MonthlyBreakdown      []MonthlySummary `json:"monthly_breakdown"`
	GeneratedAt           time.Time `json:"generated_at"`
}
// MonthlySummary for semester breakdown
type MonthlySummary struct {
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	TotalInfaq float64 `json:"total_infaq"`
}
type SemesterFinancialTurnover struct {
	TotalLoaned         float64 `json:"total_loaned"`
	TotalRepaid         float64 `json:"total_repaid"`
	AvailableForLoans   float64 `json:"available_for_loans"`
}

type BeneficiaryStatistics struct {
	TotalBorrowers      int `json:"total_borrowers"`
	ByRT1               int `json:"by_rt1"`
	ByRT2               int `json:"by_rt2"`
	ByRT3               int `json:"by_rt3"`
	CompletedBorrowers  int64 `json:"completed_borrowers"`
	ActiveBorrowers     int64 `json:"active_borrowers"`
}


type DonationCut struct {
	gorm.Model
	InfaqID uint    `json:"infaq_id" gorm:"not null;index"`
	Purpose string  `json:"purpose"`  // "dkm", "operasional", "pengurus"
	Amount  float64 `json:"amount"`
}

// Loan represents an interest-free loan (pinjaman tanpa bunga)

// LoanPayment records each weekly / partial repayment
type LoanPayment struct {
	gorm.Model
	LoanID    uint      `json:"loan_id"    gorm:"not null;index" validate:"required"`
	Amount    float64   `json:"amount"     gorm:"not null"       validate:"required,gt=0"`
	PaidAt    time.Time `json:"paid_at"`
	WeekLabel string    `json:"week_label"` // e.g. "Minggu 1", "Minggu 2"
	Notes     string    `json:"notes"`
}
