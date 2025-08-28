package controllers

import (
	"net/http"
	"time"

	"API-Books/initializer"
	"API-Books/models"

	"github.com/gin-gonic/gin"
)

func CreateInfaq(c *gin.Context) {
    var input struct {
        VillagerID     *uint    `json:"villager_id"`
        NeighborhoodID uint     `json:"neighborhood_id"`
        Amount         float64  `json:"amount"`
        DonatedAt      *time.Time `json:"donated_at"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // validate neighborhood_id (must exist)
    if input.NeighborhoodID == 0 {
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

func GetInfaqs(c *gin.Context) {
    var infaqs []models.Infaq
    if err := initializer.DB.Find(&infaqs).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, infaqs)
}


func GetInfaqByID(c *gin.Context) {
    id := c.Param("id")
    var infaq models.Infaq

    if err := initializer.DB.First(&infaq, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Infaq not found"})
        return
    }

    c.JSON(http.StatusOK, infaq)
}

func UpdateInfaq(c *gin.Context) {
    id := c.Param("id")
    var infaq models.Infaq

    if err := initializer.DB.First(&infaq, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Infaq not found"})
        return
    }

    var input struct {
        VillagerID     *uint     `json:"villager_id"`
        NeighborhoodID uint      `json:"neighborhood_id"`
        Amount         float64   `json:"amount"`
        DonatedAt      *time.Time `json:"donated_at"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if input.NeighborhoodID == 0 {
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


func DeleteInfaq(c *gin.Context) {
    id := c.Param("id")
    if err := initializer.DB.Delete(&models.Infaq{}, id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Infaq deleted successfully"})
}

