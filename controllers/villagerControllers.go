package controllers

import (
	"API-Books/initializer"
	"API-Books/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateVillager godoc
// @Summary Create a new villager
// @Description Add a villager with name, family head, and neighborhood
// @Tags villagers
// @Accept json
// @Produce json
// @Param villager body models.Villager true "Villager Data"
// @Success 200 {object} models.Villager
// @Failure 400 {object} map[string]string
// @Router /villagers [post]
func CreateVillager(ctx *gin.Context) {
    var body models.Villager

    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := initializer.DB.Create(&body).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create villager"})
        return
    }

    ctx.JSON(http.StatusOK, body)
}

// Get all villagers
func GetVillagers(ctx *gin.Context) {
    var villagers []models.Villager
    if err := initializer.DB.Find(&villagers).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch villagers"})
        return
    }

    ctx.JSON(http.StatusOK, villagers)
}

// Get villager by ID
func GetVillager(ctx *gin.Context) {
    id := ctx.Param("id")
    var villager models.Villager

    if err := initializer.DB.First(&villager, id).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Villager not found"})
        return
    }

    ctx.JSON(http.StatusOK, villager)
}

// Update villager
func UpdateVillager(ctx *gin.Context) {
    id := ctx.Param("id")
    var villager models.Villager

    if err := initializer.DB.First(&villager, id).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Villager not found"})
        return
    }

    var body models.Villager
    if err := ctx.ShouldBindJSON(&body); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    villager.Name = body.Name
    villager.FamilyHeadName = body.FamilyHeadName
    villager.NeighborhoodID = body.NeighborhoodID

    if err := initializer.DB.Save(&villager).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update villager"})
        return
    }

    ctx.JSON(http.StatusOK, villager)
}

// Delete villager
func DeleteVillager(ctx *gin.Context) {
    id := ctx.Param("id")
    if err := initializer.DB.Delete(&models.Villager{}, id).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete villager"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Villager deleted"})
}
