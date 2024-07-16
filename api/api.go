package api

import (
	"Library/api/handler"
	"Library/internal/mongodb"
	_ "Library/api/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

// Routes sets up the API routes and Swagger documentation
// @title Project: Swagger Intro
// @description This swagger UI was created in lesson
// @version 1.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func Routes(db *mongodb.LibraryMongoDb) {
	router := gin.Default()
	libraryHandler := handler.NewLibraryHandler(db)
	router.POST("/books", libraryHandler.CreateBook)
	router.GET("/books", libraryHandler.GetBooks)
    router.PUT("/books/:isbn/rent", libraryHandler.UpdateBookRentalStatus)
    router.GET("/books/author/:author", libraryHandler.GetBooksByAuthor)
    router.DELETE("/books/:isbn", libraryHandler.DeleteBookByISBN)
    router.GET("/books/aggregate/authors", libraryHandler.AggregateBooksByAuthor)



	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	

	router.Run(":7777")
}
