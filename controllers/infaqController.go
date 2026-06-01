package controllers

import (
	"net/http"
	"time"

	"API-Books/initializer"
	"API-Books/models"

	"github.com/gin-gonic/gin"
)


var input struct {
        VillagerID     *uint    `json:"villager_id"`
        NeighborhoodID string     `json:"neighborhood_id"`
        Amount         float64  `json:"amount"`
        DonatedAt      *time.Time `json:"donated_at"`
    }


type CreateInfaqRequest struct {
    VillagerID     *uint   `json:"villager_id"`
    NeighborhoodID string  `json:"neighborhood_id" binding:"required"`
    Amount         float64 `json:"amount" binding:"required,gt=0"`
    DonatedAt      *string `json:"donated_at"` // optional, ISO format
}

    
// CreateInfaq godoc
// @Summary Create a new infaq record
// @Description Create a new infaq record with auto week/month/year extraction
// @Tags Infaq
// @Accept json
// @Produce json
// @Param infaq body CreateInfaqRequest true "Infaq to create"
// @Success 201 {object} models.Infaq
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /infaqs [post]
func CreateInfaq(c *gin.Context) {
    var input CreateInfaqRequest

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Parse donated_at or default to now
    donatedAt := time.Now()
    if input.DonatedAt != nil && *input.DonatedAt != "" {
        parsed, err := time.Parse("2006-01-02", *input.DonatedAt)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid donated_at format, use YYYY-MM-DD"})
            return
        }
        donatedAt = parsed
    }

    infaq := models.Infaq{
        VillagerID:     input.VillagerID,
        NeighborhoodID: input.NeighborhoodID,
        Amount:         input.Amount,
        DonatedAt:      donatedAt,
        CollectedAt:    time.Now(),
    }

    if err := initializer.DB.Create(&infaq).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Populate derived fields for response
    infaq.WeekNumber = weekNumber(donatedAt)
    infaq.Month = int(donatedAt.Month())
    infaq.Year = donatedAt.Year()

    c.JSON(http.StatusCreated, infaq)
}

// weekNumber calculates ISO week number (1-53)
func weekNumber(t time.Time) int {
    _, week := t.ISOWeek()
    return week
}
// GetInfaqs godoc
// @Summary Get all infaq records
// @Description Get all infaq records
// @Tags Infaq
// @Accept json
// @Produce json
// @Success 200 {array} models.Infaq
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /infaqs [get]
func GetInfaqs(c *gin.Context) {
    var infaqs []models.Infaq
    if err := initializer.DB.Find(&infaqs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, infaqs)
}

// GetInfaqByID godoc
// @Summary Get an infaq record by ID
// @Description Get an infaq record by ID
// @Tags Infaq
// @Accept json
// @Produce json
// @Param id path int true "Infaq ID"
// @Success 200 {object} models.Infaq
// @Failure 404 {object} gin.H{"error": "Infaq not found"
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /infaqs/{id} [get]
func GetInfaqByID(c *gin.Context) {
    id := c.Param("id")
    var infaq models.Infaq

    if err := initializer.DB.First(&infaq, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Infaq not found"})
        return
    }

    c.JSON(http.StatusOK, infaq)
}

// UpdateInfaq godoc
// @Summary Update an infaq record
// @Description Update an infaq record by ID
// @Tags Infaq
// @Accept json
// @Produce json
// @Param id path int true "Infaq ID"
// @Param infaq body models.Infaq true "Infaq to update"
// @Success 200 {object} models.Infaq
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 404 {object} gin.H{"error": "Infaq not found"
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /infaqs/{id} [put]
func UpdateInfaq(c *gin.Context) {
    id := c.Param("id")
    var infaq models.Infaq

    if err := initializer.DB.First(&infaq, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Infaq not found"})
        return
    }

 

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.NeighborhoodID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "neighborhood_id is required"})
        return
    }

    infaq.VillagerID = input.VillagerID
    infaq.NeighborhoodID = input.NeighborhoodID
    infaq.Amount = input.Amount
    if input.DonatedAt != nil {
        infaq.DonatedAt = *input.DonatedAt
    }

    if err := initializer.DB.Save(&infaq).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, infaq)
}

// DeleteInfaq godoc
// @Summary Delete an infaq record
// @Description Delete an infaq record by ID
// @Tags Infaq
// @Param id path int true "Infaq ID"
// @Success 200 {object} gin.H{"message": "Infaq deleted successfully"
// @Failure 404 {object} gin.H{"error": "Infaq not found"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /infaqs/{id} [delete]
func DeleteInfaq(c *gin.Context) {
    id := c.Param("id")
    if err := initializer.DB.Delete(&models.Infaq{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Infaq deleted successfully"})
}

// GetInfaqsByWeek godoc
// @Summary Get infaqs by week number
// @Tags Infaq
// @Router /infaqs/weekly [get]
func GetInfaqsByWeek(c *gin.Context) {
    year := c.Query("year")
    week := c.Query("week")
    
    var infaqs []models.Infaq
    query := initializer.DB
    
    if year != "" && week != "" {
        // PostgreSQL: EXTRACT(WEEK FROM donated_at)
        query = query.Where("EXTRACT(YEAR FROM donated_at) = ? AND EXTRACT(WEEK FROM donated_at) = ?", year, week)
    }
    
    if err := query.Find(&infaqs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Populate derived fields for response
    for i := range infaqs {
        infaqs[i].WeekNumber = weekNumber(infaqs[i].DonatedAt)
        infaqs[i].Month = int(infaqs[i].DonatedAt.Month())
        infaqs[i].Year = infaqs[i].DonatedAt.Year()
    }
    
    c.JSON(http.StatusOK, infaqs)
}

// GetWeeklySummary - matches your Excel "SummaryInfaq" sheet
// GET /api/infaqs/summary/weekly?year=2025
func GetWeeklySummary(c *gin.Context) {
    year := c.DefaultQuery("year", "2025")
    
    type WeeklySummary struct {
        NeighborhoodID string  `json:"neighborhood_id"`
        Week           int     `json:"week"`
        Year           int     `json:"year"`
        TotalInfaq     float64 `json:"total_infaq"`
        ResidentCount  int64   `json:"resident_count"`
    }
    
    var results []WeeklySummary
    
    err := initializer.DB.Table("infaqs").
        Select(`
            neighborhood_id,
            EXTRACT(WEEK FROM donated_at) as week,
            EXTRACT(YEAR FROM donated_at) as year,
            SUM(amount) as total_infaq,
            COUNT(DISTINCT villager_id) as resident_count
        `).
        Where("EXTRACT(YEAR FROM donated_at) = ?", year).
        Group("neighborhood_id, week, year").
        Order("week, neighborhood_id").
        Scan(&results).Error
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, results)
}

// GetMonthlySummary - matches your Excel monthly totals
// GET /api/infaqs/summary/monthly?year=2025
func GetMonthlySummary(c *gin.Context) {
    year := c.DefaultQuery("year", "2025")
    
    type MonthlySummary struct {
        NeighborhoodID string  `json:"neighborhood_id"`
        Month          int     `json:"month"`
        Year           int     `json:"year"`
        TotalInfaq     float64 `json:"total_infaq"`
    }
    
    var results []MonthlySummary
    
    err := initializer.DB.Table("infaqs").
        Select(`
            neighborhood_id,
            EXTRACT(MONTH FROM donated_at) as month,
            EXTRACT(YEAR FROM donated_at) as year,
            SUM(amount) as total_infaq
        `).
        Where("EXTRACT(YEAR FROM donated_at) = ?", year).
        Group("neighborhood_id, month, year").
        Order("month, neighborhood_id").
        Scan(&results).Error
    
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, results)
}