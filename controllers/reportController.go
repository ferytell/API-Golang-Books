package controllers

import (
	"API-Books/initializer"
	"API-Books/models"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ============================================
// MONTHLY REPORT
// ============================================

// GenerateMonthlyReport godoc
// @Summary Generate monthly financial report
// @Description Generates comprehensive monthly report matching Excel MLAPORAN format
// @Tags reports
// @Produce json
// @Param year path int true "Year (e.g., 2025)"
// @Param month path int true "Month (1-12)"
// @Success 200 {object} models.MonthlyReport
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /reports/monthly/{year}/{month} [get]
func GenerateMonthlyReport(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil || year < 2020 || year > 2030 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	month, err := strconv.Atoi(c.Param("month"))
	if err != nil || month < 1 || month > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month (1-12)"})
		return
	}

	// Calculate period boundaries
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	report := models.MonthlyReport{
		Period:        indonesianMonth(month) + " " + strconv.Itoa(year),
		PeriodEnglish: startDate.Format("January 2006"),
		Year:          year,
		Month:         month,
		GeneratedAt:   time.Now(),
	}

	// Build each of the 4 cash books
	report.GeneralCash = buildCashBook("general", startDate, endDate, year, month)
	report.DKMCash = buildCashBook("dkm", startDate, endDate, year, month)
	report.OperationalCash = buildCashBook("operational", startDate, endDate, year, month)
	report.ManagementCash = buildCashBook("management", startDate, endDate, year, month)

	// Calculate grand total
	report.TotalAllBalance = report.GeneralCash.ClosingBalance +
		report.DKMCash.ClosingBalance +
		report.OperationalCash.ClosingBalance +
		report.ManagementCash.ClosingBalance

	// Loan statistics
	report.LoanStats = calculateMonthlyLoanStats(startDate, endDate)

	// Infaq by RT
	report.InfaqByRT = calculateInfaqByRT(startDate, endDate)

	c.JSON(http.StatusOK, report)
}

// buildCashBook constructs one of the 4 cash book summaries
func buildCashBook(bookType string, startDate, endDate time.Time, year, month int) models.CashBookSummary {
	summary := models.CashBookSummary{
		BookType: bookType,
		Entries:  []models.CashBookEntryDetail{},
	}

	// Get opening balance (previous month's closing)
	prevMonthEnd := startDate.Add(-time.Second)
	var openingBalance float64

	switch bookType {
	case "general":
		// Sum all infaq + loan repayments - loan disbursements - allocations up to prev month
		openingBalance = getGeneralOpeningBalance(prevMonthEnd)
	case "dkm":
		openingBalance = getDKMOpeningBalance(prevMonthEnd)
	case "operational":
		openingBalance = getOperationalOpeningBalance(prevMonthEnd)
	case "management":
		openingBalance = getManagementOpeningBalance(prevMonthEnd)
	}

	summary.OpeningBalance = openingBalance
	//runningBalance := openingBalance

	// Build entries based on book type
	switch bookType {
	case "general":
		entries := buildGeneralCashEntries(startDate, endDate, openingBalance)
		summary.Entries = entries
		for _, e := range entries {
			summary.TotalIncome += e.Income
			summary.TotalExpense += e.Expense
		}

	case "dkm":
		entries := buildDKMEntries(startDate, endDate, openingBalance)
		summary.Entries = entries
		for _, e := range entries {
			summary.TotalIncome += e.Income
			summary.TotalExpense += e.Expense
		}

	case "operational":
		entries := buildOperationalEntries(startDate, endDate, openingBalance)
		summary.Entries = entries
		for _, e := range entries {
			summary.TotalIncome += e.Income
			summary.TotalExpense += e.Expense
		}

	case "management":
		entries := buildManagementEntries(startDate, endDate, openingBalance)
		summary.Entries = entries
		for _, e := range entries {
			summary.TotalIncome += e.Income
			summary.TotalExpense += e.Expense
		}
	}

	// Calculate closing balance from last entry
	if len(summary.Entries) > 0 {
		summary.ClosingBalance = summary.Entries[len(summary.Entries)-1].Balance
	} else {
		summary.ClosingBalance = openingBalance
	}

	return summary
}

// General Cash Book: Infaq + Loan Repayments - Loan Disbursements - 10% allocations
func buildGeneralCashEntries(startDate, endDate time.Time, openingBalance float64) []models.CashBookEntryDetail {
	var entries []models.CashBookEntryDetail
	runningBalance := openingBalance
	entryNo := 0

	// Entry 0: Opening balance
	entries = append(entries, models.CashBookEntryDetail{
		No:          0,
		Description: "Saldo Bulan Lalu",
		Income:      0,
		Expense:     0,
		Balance:     openingBalance,
	})

	// 1. Infaq collection for the month
	var totalInfaq float64
	initializer.DB.Table("infaqs").
		Select("COALESCE(SUM(amount), 0)").
		Where("donated_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalInfaq)

	if totalInfaq > 0 {
		entryNo++
		runningBalance += totalInfaq
		entries = append(entries, models.CashBookEntryDetail{
			No:          entryNo,
			Description: fmt.Sprintf("Penerimaan Infak Bulan %s", startDate.Format("January")),
			Income:      totalInfaq,
			Expense:     0,
			Balance:     runningBalance,
		})
	}

	// 2. Loan repayments
	var totalRepayments float64
	initializer.DB.Table("loan_payments").
		Select("COALESCE(SUM(amount), 0)").
		Where("paid_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&totalRepayments)

	if totalRepayments > 0 {
		entryNo++
		runningBalance += totalRepayments
		entries = append(entries, models.CashBookEntryDetail{
			No:          entryNo,
			Description: "Penerimaan Cicilan Pinjaman",
			Income:      totalRepayments,
			Expense:     0,
			Balance:     runningBalance,
		})
	}

	// 3. Loan disbursements
	var totalDisbursed float64
	initializer.DB.Table("loans").
		Select("COALESCE(SUM(amount), 0)").
		Where("start_date BETWEEN ? AND ? AND status IN (?, ?)", startDate, endDate, "ongoing", "paid").
		Scan(&totalDisbursed)

	if totalDisbursed > 0 {
		entryNo++
		runningBalance -= totalDisbursed
		entries = append(entries, models.CashBookEntryDetail{
			No:          entryNo,
			Description: "Pengeluaran Pinjaman",
			Income:      0,
			Expense:     totalDisbursed,
			Balance:     runningBalance,
		})
	}

	// 4. DKM 10% allocation
	dkmAllocation := totalInfaq * 0.10
	if dkmAllocation > 0 {
		entryNo++
		runningBalance -= dkmAllocation
		entries = append(entries, models.CashBookEntryDetail{
			No:          entryNo,
			Description: "Infak Untuk DKM 10%/minggu*)",
			Income:      0,
			Expense:     dkmAllocation,
			Balance:     runningBalance,
		})
	}

	// 5. Operational 10% allocation
	operationalAllocation := totalInfaq * 0.10
	if operationalAllocation > 0 {
		entryNo++
		runningBalance -= operationalAllocation
		entries = append(entries, models.CashBookEntryDetail{
			No:          entryNo,
			Description: "Biaya Operasional UEM 10%/minggu*)",
			Income:      0,
			Expense:     operationalAllocation,
			Balance:     runningBalance,
		})
	}

	// 6. Management 10% allocation
	managementAllocation := totalInfaq * 0.10
	if managementAllocation > 0 {
		entryNo++
		runningBalance -= managementAllocation
		entries = append(entries, models.CashBookEntryDetail{
			No:          entryNo,
			Description: "Akomodasi Pengurus UEM 10%/minggu*)",
			Income:      0,
			Expense:     managementAllocation,
			Balance:     runningBalance,
		})
	}

	return entries
}

// DKM Cash Book: 10% infaq allocations + transfers to mosque
func buildDKMEntries(startDate, endDate time.Time, openingBalance float64) []models.CashBookEntryDetail {
	var entries []models.CashBookEntryDetail
	runningBalance := openingBalance

	// Opening balance
	entries = append(entries, models.CashBookEntryDetail{
		No:          0,
		Description: "Saldo Bulan Lalu",
		Balance:     openingBalance,
	})

	// 1. DKM allocation received
	var dkmAllocation float64
	initializer.DB.Table("infaqs").
		Select("COALESCE(SUM(amount * 0.10), 0)").
		Where("donated_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&dkmAllocation)

	if dkmAllocation > 0 {
		runningBalance += dkmAllocation
		entries = append(entries, models.CashBookEntryDetail{
			No:          1,
			Description: "Penerimaan Uang Infak Untuk DKM",
			Income:      dkmAllocation,
			Expense:     0,
			Balance:     runningBalance,
		})
	}

	// 2. Transfer to DKM chairman (if any)
	var transferred float64
	initializer.DB.Table("donation_cuts").
		Select("COALESCE(SUM(amount), 0)").
		Where("purpose = ? AND created_at BETWEEN ? AND ?", "dkm_transfer", startDate, endDate).
		Scan(&transferred)

	if transferred > 0 {
		runningBalance -= transferred
		entries = append(entries, models.CashBookEntryDetail{
			No:          2,
			Description: "Diserah terimakan kepada ketua DKM",
			Income:      0,
			Expense:     transferred,
			Balance:     runningBalance,
		})
	}

	return entries
}

// Operational Cash Book: 10% for operations
func buildOperationalEntries(startDate, endDate time.Time, openingBalance float64) []models.CashBookEntryDetail {
	var entries []models.CashBookEntryDetail
	runningBalance := openingBalance

	entries = append(entries, models.CashBookEntryDetail{
		No:          0,
		Description: "Saldo Bulan Lalu",
		Balance:     openingBalance,
	})

	// 1. Operational allocation received
	var operationalAllocation float64
	initializer.DB.Table("infaqs").
		Select("COALESCE(SUM(amount * 0.10), 0)").
		Where("donated_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&operationalAllocation)

	if operationalAllocation > 0 {
		runningBalance += operationalAllocation
		entries = append(entries, models.CashBookEntryDetail{
			No:          1,
			Description: fmt.Sprintf("Pemasukan Bulan %s", startDate.Format("January")),
			Income:      operationalAllocation,
			Expense:     0,
			Balance:     runningBalance,
		})
	}

	// 2. Operational expenses
	var expenses float64
	initializer.DB.Table("donation_cuts").
		Select("COALESCE(SUM(amount), 0)").
		Where("purpose = ? AND created_at BETWEEN ? AND ?", "operational", startDate, endDate).
		Scan(&expenses)

	if expenses > 0 {
		runningBalance -= expenses
		entries = append(entries, models.CashBookEntryDetail{
			No:          2,
			Description: "Pembelian Keperluan Operasional & Usaha UEM",
			Income:      0,
			Expense:     expenses,
			Balance:     runningBalance,
		})
	}

	return entries
}

// Management Cash Book: 10% for management transport/accommodation
func buildManagementEntries(startDate, endDate time.Time, openingBalance float64) []models.CashBookEntryDetail {
	var entries []models.CashBookEntryDetail
	runningBalance := openingBalance

	entries = append(entries, models.CashBookEntryDetail{
		No:          0,
		Description: "Saldo Bulan Lalu",
		Balance:     openingBalance,
	})

	// 1. Management allocation received
	var managementAllocation float64
	initializer.DB.Table("infaqs").
		Select("COALESCE(SUM(amount * 0.10), 0)").
		Where("donated_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&managementAllocation)

	if managementAllocation > 0 {
		runningBalance += managementAllocation
		entries = append(entries, models.CashBookEntryDetail{
			No:          1,
			Description: "Penerimaan Uang Akomodasi Pengurus UEM",
			Income:      managementAllocation,
			Expense:     0,
			Balance:     runningBalance,
		})
	}

	// 2. Management expenses (transport, etc.)
	var expenses float64
	initializer.DB.Table("donation_cuts").
		Select("COALESCE(SUM(amount), 0)").
		Where("purpose = ? AND created_at BETWEEN ? AND ?", "management", startDate, endDate).
		Scan(&expenses)

	if expenses > 0 {
		runningBalance -= expenses
		entries = append(entries, models.CashBookEntryDetail{
			No:          2,
			Description: "Pembayaran Akomodasi Pengurus UEM",
			Income:      0,
			Expense:     expenses,
			Balance:     runningBalance,
		})
	}

	return entries
}

// ============================================
// SEMESTER REPORT
// ============================================

// GenerateSemesterReport godoc
// @Summary Generate semester financial report
// @Description Generates 6-month semester report (Jan-Jun or Jul-Dec)
// @Tags reports
// @Produce json
// @Param year path int true "Year"
// @Param semester path int true "Semester (1 or 2)"
// @Success 200 {object} models.SemesterReport
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /reports/semester/{year}/{semester} [get]
func GenerateSemesterReport(c *gin.Context) {
	year, err := strconv.Atoi(c.Param("year"))
	if err != nil || year < 2020 || year > 2030 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}

	semester, err := strconv.Atoi(c.Param("semester"))
	if err != nil || (semester != 1 && semester != 2) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "semester must be 1 (Jan-Jun) or 2 (Jul-Dec)"})
		return
	}

	var startMonth, endMonth int
	var periodName string

	if semester == 1 {
		startMonth = 1
		endMonth = 6
		periodName = fmt.Sprintf("Januari - Juni %d", year)
	} else {
		startMonth = 7
		endMonth = 12
		periodName = fmt.Sprintf("Juli - Desember %d", year)
	}

	startDate := time.Date(year, time.Month(startMonth), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, time.Month(endMonth+1), 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

	report := models.SemesterReport{
		Period:     periodName,
		Semester:   semester,
		Year:       year,
		GeneratedAt: time.Now(),
	}

	// I. Total Infaq Collected
	initializer.DB.Table("infaqs").
		Select("COALESCE(SUM(amount), 0)").
		Where("donated_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&report.TotalInfaqCollected)

	// II. Total Loan Repayments
	initializer.DB.Table("loan_payments").
		Select("COALESCE(SUM(amount), 0)").
		Where("paid_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&report.TotalLoanRepayments)

	// III. Total Loan Disbursed
	initializer.DB.Table("loans").
		Select("COALESCE(SUM(amount), 0)").
		Where("start_date BETWEEN ? AND ?", startDate, endDate).
		Scan(&report.TotalLoanDisbursed)

	// IV. DKM Transfers (10% of infaq)
	report.TotalDKMTransferred = report.TotalInfaqCollected * 0.10

	// V. Operational Used (10%)
	report.TotalOperationalUsed = report.TotalInfaqCollected * 0.10

	// VI. Management Used (10%)
	report.TotalManagementUsed = report.TotalInfaqCollected * 0.10

	// Financial Turnover
	report.FinancialTurnover = calculateSemesterFinancialTurnover(startDate, endDate)

	// Beneficiary Statistics
	report.BeneficiaryStats = calculateBeneficiaryStats(startDate, endDate)

	// Monthly breakdown
	//report.MonthlyBreakdown = buildMonthlyBreakdown(year, startMonth, endMonth)

	c.JSON(http.StatusOK, report)
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func calculateMonthlyLoanStats(startDate, endDate time.Time) models.MonthlyLoanStats {
	var stats models.MonthlyLoanStats

	// New loans this month
	initializer.DB.Model(&models.Loan{}).
		Where("start_date BETWEEN ? AND ?", startDate, endDate).
		Count(&stats.NewLoans)

	// Total disbursed
	initializer.DB.Table("loans").
		Select("COALESCE(SUM(amount), 0)").
		Where("start_date BETWEEN ? AND ?", startDate, endDate).
		Scan(&stats.TotalDisbursed)

	// Total repaid
	initializer.DB.Table("loan_payments").
		Select("COALESCE(SUM(amount), 0)").
		Where("paid_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&stats.TotalRepaid)

	// Active loans
	initializer.DB.Model(&models.Loan{}).
		Where("status = ?", "ongoing").
		Count(&stats.ActiveLoans)

	// Completed loans this month
	initializer.DB.Model(&models.Loan{}).
		Where("status = ? AND actual_end_date BETWEEN ? AND ?", "paid", startDate, endDate).
		Count(&stats.CompletedLoans)

	return stats
}

func calculateInfaqByRT(startDate, endDate time.Time) []models.RTInfaqSummary {
	var results []models.RTInfaqSummary

	initializer.DB.Table("infaqs").
		Select(`
			neighborhood_id,
			COALESCE(SUM(amount), 0) as total_infaq,
			COUNT(DISTINCT villager_id) as resident_count
		`).
		Where("donated_at BETWEEN ? AND ?", startDate, endDate).
		Group("neighborhood_id").
		Order("neighborhood_id").
		Scan(&results)

	return results
}

func calculateSemesterFinancialTurnover(startDate, endDate time.Time) models.SemesterFinancialTurnover {
	var ft models.SemesterFinancialTurnover

	// Total ever loaned (cumulative)
	initializer.DB.Table("loans").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&ft.TotalLoaned)

	// Total ever repaid (cumulative)
	initializer.DB.Table("loan_payments").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&ft.TotalRepaid)

	// Available for loans = current general cash balance
	//var generalBalance float64
	// This would need to query the actual balance from cash book entries
	// For now, approximate:
	ft.AvailableForLoans = ft.TotalLoaned - ft.TotalRepaid

	return ft
}

func calculateBeneficiaryStats(startDate, endDate time.Time) models.BeneficiaryStatistics {
	var stats models.BeneficiaryStatistics

	// Total borrowers
	initializer.DB.Model(&models.Loan{}).
		Select("COUNT(DISTINCT villager_id)").
		Scan(&stats.TotalBorrowers)

	// By RT
	initializer.DB.Table("loans").
		Select("COUNT(DISTINCT loans.villager_id)").
		Joins("JOIN villagers ON villagers.id = loans.villager_id").
		Where("villagers.neighborhood_id = ?", "01").
		Scan(&stats.ByRT1)

	initializer.DB.Table("loans").
		Select("COUNT(DISTINCT loans.villager_id)").
		Joins("JOIN villagers ON villagers.id = loans.villager_id").
		Where("villagers.neighborhood_id = ?", "02").
		Scan(&stats.ByRT2)

	initializer.DB.Table("loans").
		Select("COUNT(DISTINCT loans.villager_id)").
		Joins("JOIN villagers ON villagers.id = loans.villager_id").
		Where("villagers.neighborhood_id = ?", "03").
		Scan(&stats.ByRT3)

	// Completed
	initializer.DB.Model(&models.Loan{}).
		Where("status = ?", "paid").
		Count(&stats.CompletedBorrowers)

	// Active
	initializer.DB.Model(&models.Loan{}).
		Where("status = ?", "ongoing").
		Count(&stats.ActiveBorrowers)

	return stats
}

func buildMonthlyBreakdown(year, startMonth, endMonth int) []models.MonthlySummary {
	var breakdown []models.MonthlySummary

	for m := startMonth; m <= endMonth; m++ {
		startDate := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

		var summary models.MonthlySummary
		summary.Month = m
		summary.Year = year

		initializer.DB.Table("infaqs").
			Select("COALESCE(SUM(amount), 0) as total_infaq").
			Where("donated_at BETWEEN ? AND ?", startDate, endDate).
			Scan(&summary.TotalInfaq)

		breakdown = append(breakdown, summary)
	}

	return breakdown
}

// Opening balance helpers (simplified - in production, query actual closing balance from previous period)
func getGeneralOpeningBalance(asOf time.Time) float64 {
	var balance float64
	// Sum all income - all expenses up to this date
	// This is a simplified calculation
	initializer.DB.Raw(`
		SELECT COALESCE(
			(SELECT SUM(amount) FROM infaqs WHERE donated_at <= ?) +
			(SELECT SUM(amount) FROM loan_payments WHERE paid_at <= ?) -
			(SELECT SUM(amount) FROM loans WHERE start_date <= ?) -
			(SELECT SUM(amount * 0.30) FROM infaqs WHERE donated_at <= ?),
			0
		)
	`, asOf, asOf, asOf, asOf).Scan(&balance)
	return balance
}

func getDKMOpeningBalance(asOf time.Time) float64 {
	var balance float64
	initializer.DB.Raw(`
		SELECT COALESCE(
			(SELECT SUM(amount * 0.10) FROM infaqs WHERE donated_at <= ?) -
			(SELECT SUM(amount) FROM donation_cuts WHERE purpose = 'dkm_transfer' AND created_at <= ?),
			0
		)
	`, asOf, asOf).Scan(&balance)
	return balance
}

func getOperationalOpeningBalance(asOf time.Time) float64 {
	var balance float64
	initializer.DB.Raw(`
		SELECT COALESCE(
			(SELECT SUM(amount * 0.10) FROM infaqs WHERE donated_at <= ?) -
			(SELECT SUM(amount) FROM donation_cuts WHERE purpose = 'operational' AND created_at <= ?),
			0
		)
	`, asOf, asOf).Scan(&balance)
	return balance
}

func getManagementOpeningBalance(asOf time.Time) float64 {
	var balance float64
	initializer.DB.Raw(`
		SELECT COALESCE(
			(SELECT SUM(amount * 0.10) FROM infaqs WHERE donated_at <= ?) -
			(SELECT SUM(amount) FROM donation_cuts WHERE purpose = 'management' AND created_at <= ?),
			0
		)
	`, asOf, asOf).Scan(&balance)
	return balance
}

func indonesianMonth(m int) string {
	months := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}
	if m >= 1 && m <= 12 {
		return months[m]
	}
	return ""
}
