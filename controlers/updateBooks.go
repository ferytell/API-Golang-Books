package controlers

import (
	"API-Books/initializer"
	"API-Books/models"

	"github.com/gin-gonic/gin"
)

func UpdateBook(ctx *gin.Context) {

	// Get id
	id := ctx.Param("id")

	// Get Data off requested Body

	var body struct {
		Title       string
		Author      string
		Description string
		Body        string
	}

	ctx.Bind(&body)

	// Find the post were updating

	var post models.Post
	initializer.DB.First(&post, id)

	// Update the Database

	initializer.DB.Model(&post).Updates(models.Post{
		Title:       body.Title,
		Author:      body.Author,
		Description: body.Description,
		Body:        body.Body,
	})

	// Response

	ctx.JSON(200, gin.H{
		"post": post,
	})

	// bookId := ctx.Param("bookId")
	// condition := false
	// var updateBook Book

	// if err := ctx.ShouldBindJSON(&updateBook); err != nil {
	// 	ctx.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }

	// for i, book := range BookDatas {
	// 	if bookId == book.BookId {
	// 		condition = true
	// 		BookDatas[i] = updateBook
	// 		BookDatas[i].BookId = bookId
	// 		break
	// 	}
	// }

	// if !condition {
	// 	ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
	// 		"error_status":  "Data Not Found",
	// 		"error_message": fmt.Sprintf("Book with id %v not found", bookId),
	// 	})
	// 	return
	// }

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"message": fmt.Sprintf("Book with id %v has been updated", bookId),
	// })
}
