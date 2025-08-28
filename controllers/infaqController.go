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
    
// CreateInfaq godoc
// @Summary Create a new infaq record
// @Description Create a new infaq record
// @Tags Infaq
// @Accept json
// @Produce json
// @Param infaq body models.Infaq true "Infaq to create"
// @Success 201 {object} models.Infaq
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /infaqs [post]
func CreateInfaq(c *gin.Context) {
    

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // validate neighborhood_id (must exist)
    if input.NeighborhoodID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "neighborhood_id is required"})
        return
    }

    donatedAt := time.Now()
    if input.DonatedAt != nil {
        donatedAt = *input.DonatedAt
    }

    infaq := models.Infaq{
        VillagerID:     input.VillagerID,
        NeighborhoodID: input.NeighborhoodID,
        Amount:         input.Amount,
        DonatedAt:      donatedAt,
    }

    if err := initializer.DB.Create(&infaq).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, infaq)
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

