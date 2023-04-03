package controlers

import (
	"API-Books/initializer"
	"API-Books/models"

	"github.com/gin-gonic/gin"
)

func CreateBook(ctx *gin.Context) {

	// Get data off req body

	var body struct {
		Title       string
		Author      string
		Description string
		Body        string
	}

	ctx.Bind(&body)
	// Create Post Books
	post := models.Post{
		Title:       body.Title,
		Author:      body.Author,
		Description: body.Description,
		Body:        body.Body,
	}

	result := initializer.DB.Create(&post)

	if result.Error != nil {
		ctx.Status(400)
		return
	}

	// Return It

	ctx.JSON(200, gin.H{
		"post": post,
	})

	// var newBook Book

	// if err := ctx.ShouldBindJSON(&newBook); err != nil {
	// 	ctx.AbortWithError(http.StatusBadRequest, err)
	// 	return
	// }

	// newBook.BookId = fmt.Sprintf("c%d", len(BookDatas)+1)
	// BookDatas = append(BookDatas, newBook)

	// ctx.JSON(http.StatusCreated, gin.H{
	// 	"book": newBook,
	// })
}
