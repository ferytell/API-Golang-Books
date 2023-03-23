package controlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetBook(ctx *gin.Context) {
	bookId := ctx.Param("bookId")
	condition := false
	var bookData Book

	for i, book := range BookDatas {

		if bookId == book.BookId {
			condition = true
			bookData = BookDatas[i]
			break
		}
	}

	if !condition {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Data Not Found",
			"error_message": fmt.Sprintf("book with id %v not Found", bookId),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"book": bookData,
	})

}

func GetAllBooks(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"books": BookDatas,
	})
}
