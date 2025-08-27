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

// GetVillagers godoc
// @Summary Get all villagers
// @Description Retrieve a list of all villagers
// @Tags villagers
// @Produce json
// @Success 200 {array} models.Villager
// @Failure 500 {object} map[string]string
// @Router /villagers [get]
func GetVillagers(ctx *gin.Context) {
    var villagers []models.Villager
    if err := initializer.DB.Find(&villagers).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch villagers"})
        return
    }

    ctx.JSON(http.StatusOK, villagers)
}

// GetVillagerById godoc
// @Summary Get a villager by ID
// @Description Retrieve a single villager by their ID
// @Tags villagers
// @Produce json
// @Param id path int true "Villager ID"
// @Success 200 {object} models.Villager
// @Failure 404 {object} map[string]string
// @Router /villagers/{id} [get]
func GetVillager(ctx *gin.Context) {
    id := ctx.Param("id")
    var villager models.Villager

    if err := initializer.DB.First(&villager, id).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Villager not found"})
        return
    }

    ctx.JSON(http.StatusOK, villager)
}

// UpdateVillager godoc
// @Summary Update a villager
// @Description Update villager details by ID
// @Tags villagers
// @Accept json
// @Produce json
// @Param id path int true "Villager ID"
// @Param villager body models.Villager true "Updated Villager Data"
// @Success 200 {object} models.Villager
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /villagers/{id} [put]
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

// DeleteVillager godoc
// @Summary Delete a villager
// @Description Remove a villager by ID
// @Tags villagers
// @Param id path int true "Villager ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /villagers/{id} [delete]
func DeleteVillager(ctx *gin.Context) {
    id := ctx.Param("id")
    if err := initializer.DB.Delete(&models.Villager{}, id).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete villager"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"message": "Villager deleted"})
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

func GetVillagersWithLoans(ctx *gin.Context) {
    var villagers []models.Villager
    if err := initializer.DB.Preload("Loans").Find(&villagers).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch"})
        return
    }
    ctx.JSON(http.StatusOK, villagers)
}
