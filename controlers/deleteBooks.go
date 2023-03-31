package controlers

import (
	"API-Books/initializer"
	"API-Books/models"

	"github.com/gin-gonic/gin"
)

func DeleteBook(ctx *gin.Context) {
	// Get id
	id := ctx.Param("id")

	// Delete Data
	initializer.DB.Delete(&models.Post{}, id)

	// Response
	ctx.Status(200)

	// 	bookId := ctx.Param("bookId")
	// 	condition := false
	// 	var bookIndex int

	// 	for i, book := range BookDatas {
	// 		if bookId == book.BookId {
	// 			condition = true
	// 			bookIndex = i
	// 			break
	// 		}
	// 	}

	// 	if !condition {
	// 		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
	// 			"error_status":  "Data Noot Found",
	// 			"error_message": fmt.Sprintf("book with id %v not found", bookId),
	// 		})
	// 		return
	// 	}

	// 	copy(BookDatas[bookIndex:], BookDatas[bookIndex+1:])
	// 	BookDatas[len(BookDatas)-1] = Book{}
	// 	BookDatas = BookDatas[:len(BookDatas)-1]

	// 	ctx.JSON(http.StatusOK, gin.H{
	// 		"message": fmt.Sprintf("book with id %v has been deleted", bookId),
	// 	})
	// }

	// func TestDel(ctx *gin.Context) {
	// 	bookId := ctx.Param("bookId")
	// 	ctx.JSON(200, gin.H{
	// 		"message": fmt.Sprintf("book with id %v has been deleted", bookId),
	// 	})
}
