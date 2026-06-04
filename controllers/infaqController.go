package controllers

import (
	"net/http"
	"time"

	"API-Books/initializer"
	"API-Books/models"

	"github.com/gin-gonic/gin"
)

// infaqInput is the request body for create/update — defined locally (not as a package var)
type infaqInput struct {
	VillagerID     *uint      `json:"villager_id"`
	NeighborhoodID string     `json:"neighborhood_id" binding:"required"`
	Amount         float64    `json:"amount"          binding:"required,gt=0"`
	DonatedAt      *time.Time `json:"donated_at"`
	Notes          string     `json:"notes"`
}

// CreateInfaq godoc
// @Summary Create a new infaq record
// @Tags Infaq
// @Accept json
// @Produce json
// @Param infaq body infaqInput true "Infaq data"
// @Success 201 {object} models.Infaq
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/infaqs [post]
func CreateInfaq(c *gin.Context) {
	var body infaqInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	donatedAt := time.Now()
	if body.DonatedAt != nil {
		donatedAt = *body.DonatedAt 
	}

	infaq := models.Infaq{
		VillagerID:     body.VillagerID,
		NeighborhoodID: body.NeighborhoodID,
		Amount:         body.Amount,
		DonatedAt:      donatedAt,
		Notes:          body.Notes,
		CollectedAt:    time.Now(),
	}

	if err := initializer.DB.Create(&infaq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, infaq)
}

// weekNumber calculates ISO week number (1-53)
func weekNumber(t time.Time) int {
    _, week := t.ISOWeek()
    return week
}
// GetInfaqs godoc
// @Summary Get all infaq records
// @Tags Infaq
// @Produce json
// @Param neighborhood_id query string false "Filter by RT (e.g. 01)"
// @Param year            query int    false "Filter by year"
// @Param month           query int    false "Filter by month (1-12)"
// @Success 200 {array} models.Infaq
// @Failure 500 {object} map[string]string
// @Router /api/infaqs [get]
func GetInfaqs(c *gin.Context) {
	var infaqs []models.Infaq
	q := initializer.DB.Model(&models.Infaq{})

	if rt := c.Query("neighborhood_id"); rt != "" {
		q = q.Where("neighborhood_id = ?", rt)
	}
	if year := c.Query("year"); year != "" {
		q = q.Where("EXTRACT(YEAR FROM donated_at) = ?", year)
	}
	if month := c.Query("month"); month != "" {
		q = q.Where("EXTRACT(MONTH FROM donated_at) = ?", month)
	}

	if err := q.Order("donated_at desc").Find(&infaqs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, infaqs)
}

// GetInfaqByID godoc
// @Summary Get infaq by ID
// @Tags Infaq
// @Produce json
// @Param id path int true "Infaq ID"
// @Success 200 {object} models.Infaq
// @Failure 404 {object} map[string]string
// @Router /api/infaqs/{id} [get]
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
// @Tags Infaq
// @Accept json
// @Produce json
// @Param id   path int        true "Infaq ID"
// @Param body body infaqInput true "Updated infaq"
// @Success 200 {object} models.Infaq
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/infaqs/{id} [put]
func UpdateInfaq(c *gin.Context) {
	id := c.Param("id")
	var infaq models.Infaq
	if err := initializer.DB.First(&infaq, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Infaq not found"})
		return
	}

	var body infaqInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	infaq.VillagerID = body.VillagerID
	infaq.NeighborhoodID = body.NeighborhoodID
	infaq.Amount = body.Amount
	infaq.Notes = body.Notes
	if body.DonatedAt != nil {
		infaq.DonatedAt = *body.DonatedAt
	}

	if err := initializer.DB.Save(&infaq).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, infaq)
}

// DeleteInfaq godoc
// @Summary Delete an infaq record
// @Tags Infaq
// @Param id path int true "Infaq ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/infaqs/{id} [delete]
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