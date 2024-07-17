package api

import (
	"Library/api/handler"
	"Library/internal/mongodb"

	"github.com/gin-gonic/gin"

)

func Routes(db *mongodb.LibraryMongoDb) {
	router := gin.Default()
	libraryHandler := handler.NewLibraryHandler(db)
	router.POST("/books", libraryHandler.CreateBook)
	router.GET("/books", libraryHandler.GetBooks)
    router.PUT("/books/:isbn/rent", libraryHandler.UpdateBookRentalStatus)
    router.GET("/books/author/:author", libraryHandler.GetBooksByAuthor)
    router.DELETE("/books/:isbn", libraryHandler.DeleteBookByISBN)
    router.GET("/books/aggregate/authors", libraryHandler.AggregateBooksByAuthor)

	router.Run(":7777")
}
