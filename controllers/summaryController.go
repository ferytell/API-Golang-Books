package controllers

import (
	"API-Books/initializer"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInfaqSummaryByWeek(ctx *gin.Context) {
    type Result struct {
        NeighborhoodID string
        Week           int
        Year           int
        TotalInfaq     float64
    }

    var results []Result
    err := initializer.DB.Table("infaqs").
        Select(`
            villagers.neighborhood_id,
            EXTRACT(WEEK FROM infaqs.collected_at) as week,
            EXTRACT(YEAR FROM infaqs.collected_at) as year,
            SUM(infaqs.amount) as total_infaq`).
        Joins("LEFT JOIN villagers ON villagers.id = infaqs.villager_id").
        Group("villagers.neighborhood_id, week, year").
        Order("year, week").
        Scan(&results).Error

    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
        return
    }

    ctx.JSON(http.StatusOK, results)
}



func GetInfaqSummaryByMonth(ctx *gin.Context) {
    type Result struct {
        NeighborhoodID string
        Month          int
        Year           int
        TotalInfaq     float64
    }

    var results []Result
    err := initializer.DB.Table("infaqs").
        Select(`
            villagers.neighborhood_id,
            EXTRACT(MONTH FROM infaqs.collected_at) as month,
            EXTRACT(YEAR FROM infaqs.collected_at) as year,
            SUM(infaqs.amount) as total_infaq`).
        Joins("LEFT JOIN villagers ON villagers.id = infaqs.villager_id").
        Group("villagers.neighborhood_id, month, year").
        Order("year, month").
        Scan(&results).Error

    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch summary"})
        return
    }

    ctx.JSON(http.StatusOK, results)
}


func GetFundAllocation(ctx *gin.Context) {
    type Result struct {
        Purpose string
        Total   float64
    }

    var results []Result
    err := initializer.DB.Table("donation_cuts").
        Select("purpose, SUM(amount) as total").
        Group("purpose").
        Scan(&results).Error

    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch fund allocation"})
        return
    }

    ctx.JSON(http.StatusOK, results)
}
