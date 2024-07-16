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

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book record in the library
// @Tags Book
// @Accept json
// @Produce json
// @Param body body models.Book true "Book object to be created"
// @Security BearerAuth
// @Success 201 {object} gin.H{"inserted": string} "Successfully created"
// @Failure 400 {object} gin.H{"error": string} "Bad request"
// @Failure 500 {object} gin.H{"error": string} "Internal server error"
// @Router /books [post]
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

// GetBooks godoc
// @Summary Get all books
// @Description Get all books currently in the library
// @Tags Book
// @Accept json
// @Produce json
// @Success 200 {array} models.Book "List of books"
// @Failure 500 {object} gin.H{"error": string} "Internal server error"
// @Router /books [get]
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

// UpdateBookRentalStatus godoc
// @Summary Update book rental status
// @Description Update the rental status (rented or available) of a specific book
// @Tags Book
// @Accept json
// @Produce json
// @Param isbn path string true "ISBN of the book to update"
// @Param body body struct{ IsRented bool `json:"is_rented"` } true "Rental status"
// @Security BearerAuth
// @Success 200 {object} gin.H{"message": string} "Successfully updated"
// @Failure 400 {object} gin.H{"error": string} "Bad request"
// @Failure 500 {object} gin.H{"error": string} "Internal server error"
// @Router /books/{isbn} [put]
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

// GetBooksByAuthor godoc
// @Summary Get books by author
// @Description Get all books written by a specific author
// @Tags Book
// @Accept json
// @Produce json
// @Param author path string true "Author's name"
// @Success 200 {array} models.Book "List of books by the author"
// @Failure 500 {object} gin.H{"error": string} "Internal server error"
// @Router /books/byauthor/{author} [get]
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

// DeleteBookByISBN godoc
// @Summary Delete a book by ISBN
// @Description Delete a book from the library by its ISBN
// @Tags Book
// @Accept json
// @Produce json
// @Param isbn path string true "ISBN of the book to delete"
// @Security BearerAuth
// @Success 200 {object} gin.H{"message": string} "Successfully deleted"
// @Failure 500 {object} gin.H{"error": string} "Internal server error"
// @Router /books/{isbn} [delete]
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

// AggregateBooksByAuthor godoc
// @Summary Aggregate books by author
// @Description Aggregate and get book count by each author
// @Tags Book
// @Accept json
// @Produce json
// @Success 200 {object} gin.H{"author": string, "count": int} "Aggregated results"
// @Failure 500 {object} gin.H{"error": string} "Internal server error"
// @Router /books/aggregate/byauthor [get]
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
