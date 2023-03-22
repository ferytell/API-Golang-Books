package controlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	BookId   string `json:"book_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
}

var BookDatas = []Book{}

func CreateBook(ctx *gin.Context) {
	var newBook Book

	if err := ctx.ShouldBindJSON(&newBook); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	newBook.BookId = fmt.Sprintf("c%d", len(BookDatas)+1)
	BookDatas = append(BookDatas, newBook)

	ctx.JSON(http.StatusCreated, gin.H{
		"book": newBook,
	})
}

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
			//BookDatas[i].BookId = bookId
			break
		}
	}

	if condition {
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

func Hellow(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello bro",
	})
}

func Home(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Homeee",
	})
}

func GetAllBooks(c *gin.Context) {

}
