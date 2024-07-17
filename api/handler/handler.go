package handler

import (
	"Library/internal/models"
	"Library/internal/mongodb"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LibraryHandler struct {
	db *mongodb.LibraryMongoDb
}

func NewLibraryHandler(db *mongodb.LibraryMongoDb) *LibraryHandler {
	return &LibraryHandler{db: db}
}

func (l *LibraryHandler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := l.db.InsertBook(ctx, &book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"inserted": result})
}


func (l *LibraryHandler) GetBooks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	books, err := l.db.GetAllBooks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}


func (l *LibraryHandler) UpdateBookRentalStatus(c *gin.Context) {
	isbn := c.Param("isbn")
	var status struct {
		IsRented bool `json:"is_rented"`
	}
	if err := c.BindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := l.db.UpdateBooks(ctx, isbn, status.IsRented)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated the rental status of the book"})
}


func (l *LibraryHandler) GetBooksByAuthor(c *gin.Context) {
	author := c.Param("author")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	books, err := l.db.GetAllbyAuthor(ctx, author)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

func (l *LibraryHandler) DeleteBookByISBN(c *gin.Context) {
	isbn := c.Param("isbn")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := l.db.DeleteBookByISBN(ctx, isbn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted the book"})
}


func (l *LibraryHandler) AggregateBooksByAuthor(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	results, err := l.db.AggregateBooksByAuthor(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
