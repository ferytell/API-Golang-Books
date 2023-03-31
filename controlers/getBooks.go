package controlers

import (
	"API-Books/initializer"
	"API-Books/models"

	"github.com/gin-gonic/gin"
)

func GetBook(ctx *gin.Context) {

	id := ctx.Param("id")
	var post models.Post

	initializer.DB.First(&post, id)

	ctx.JSON(200, gin.H{
		"post": post,
	})
	// bookId := ctx.Param("bookId")
	// condition := false
	// var bookData Book

	// for i, book := range BookDatas {

	// 	if bookId == book.BookId {
	// 		condition = true
	// 		bookData = BookDatas[i]
	// 		break
	// 	}
	// }

	// if !condition {
	// 	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
	// 		"error_status":  "Data Not Found",
	// 		"error_message": fmt.Sprintf("book with id %v not Found", bookId),
	// 	})
	// 	return
	// }

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"book": bookData,
	// })

}

func GetAllBooks(ctx *gin.Context) {
	// Get The Books

	var posts []models.Post

	initializer.DB.Find(&posts)

	ctx.JSON(200, gin.H{
		"posts": posts,
	})
}
