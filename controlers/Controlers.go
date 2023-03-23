package controlers

import (
	"github.com/gin-gonic/gin"
)

var BookDatas = []Book{}

type Book struct {
	BookId      string `json:"book_id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"desc"`
}

func Hellow(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Hello bro, tester aja gw mah",
	})
}

func Home(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "Homeee",
	})
}
