package controlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateBook(ctx *gin.Context) {
	bookId := ctx.Param("bookId")
	condition := false
	var updateBook Book

	if err := ctx.ShouldBindJSON(&updateBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	for i, book := range BookDatas {
		if bookId == book.BookId {
			condition = true
			BookDatas[i] = updateBook
			BookDatas[i].BookId = bookId
			break
		}
	}

	if !condition {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error_status":  "Data Not Found",
			"error_message": fmt.Sprintf("Book with id %v not found", bookId),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Book with id %v has been updated", bookId),
	})
}
