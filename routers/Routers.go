package routers

import (
	"API-Books/controllers"
	"API-Books/middleware"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.New() // Use gin.New(), not gin.Default() for more control
    router.Use(gin.Recovery())


	 // CORS FIRST
    router.Use(cors.New(cors.Config{
        AllowOriginFunc: func(origin string) bool {
            return origin == "http://localhost:3000" ||
                   strings.HasPrefix(origin, "https://ferytell.site") ||
                   strings.HasPrefix(origin, "https://testing.ferytell.site")
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Then logger (optional)
    router.Use(gin.Logger())

	//router.Use(cors.New(config))

	router.GET("/api/ping", controllers.Hellow)
	// User Auth Routes
	router.POST("/api/signup", controllers.SignUp)
	router.POST("/api/login", controllers.Login)
	router.GET("/api/validate", middleware.RequireAuth, controllers.Validate)
	router.POST("/api/logout", controllers.Logout)

	// Books Routes
	router.GET("/api/books", controllers.GetAllBooks)
	router.GET("/api/books/:id", controllers.GetBook)
	router.POST("/api/books", controllers.CreateBook)
	router.PUT("/api/books/:id", controllers.UpdateBook)
	router.DELETE("/api/books/:id", controllers.DeleteBook)

	//Villager Routes
	router.POST("/api/villagers", controllers.CreateVillager)
	router.GET("/api/villagers", controllers.GetVillagers)
	router.GET("/api/villagers/:id", controllers.GetVillager)
	router.PUT("/api/villagers/:id", controllers.UpdateVillager)
	router.GET("/api/villagers/:id/loans", controllers.GetVillagerLoans)
	router.DELETE("/api/villagers/:id", controllers.DeleteVillager)

	//Infaq Routes
	router.POST("/api/infaqs", controllers.CreateInfaq)
	router.GET("/api/infaqs", controllers.GetInfaqs)
	router.GET("/api/infaqs/:id", controllers.GetInfaqByID)
	router.PUT("/api/infaqs/:id", controllers.UpdateInfaq)
	router.DELETE("/api/infaqs/:id", controllers.DeleteInfaq)
	

	// Summary Routes
	router.GET("/api/summary/infaq/week", controllers.GetInfaqSummaryByWeek)
	router.GET("/api/summary/infaq/month", controllers.GetInfaqSummaryByMonth)
	router.GET("/summary/fund-allocation", controllers.GetFundAllocation)

	// Loan Routes
	router.POST("/api/loans", controllers.RequestLoan)
	router.POST("/api/loans/batch", controllers.PaymentLoan)
	router.GET("/api/loans", controllers.GetAllLoans)
	
	// Payment Routes
	router.POST("/api/loan_payments", controllers.CreatePayment)
	router.PUT("/api/loan_payments/:id", controllers.UpdatePayment)
	router.GET("/api/loan_payments", controllers.GetAllPayments)
	router.GET("/api/loan_payments/:id", controllers.GetVillagersWithLoans)


	return router
}
