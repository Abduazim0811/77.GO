package mongodb

import (
	"Library/internal/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LibraryMongoDb struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func MongoDb(url, dbname, collectionName string) (*LibraryMongoDb, error) {
	clientOptions := options.Client().ApplyURI(url)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("ping error: %v", err)
	}

	collection := client.Database(dbname).Collection(collectionName)
	return &LibraryMongoDb{Client: client, Collection: collection}, nil
}

func (l *LibraryMongoDb) InsertBook(ctx context.Context, book *models.Book) (*mongo.InsertOneResult, error) {
	return l.Collection.InsertOne(ctx, book)
}
func (l *LibraryMongoDb) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	cursor, err := l.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var books []models.Book
	for cursor.Next(ctx) {
		var book models.Book
		cursor.Decode(&book)
		books = append(books, book)
	}

	return books, nil
}

func (l *LibraryMongoDb) UpdateBooks(ctx context.Context, isbn string, isRented bool) error {
	filter := bson.M{"isbn": isbn}
	update := bson.M{"$set": bson.M{"is_rented": isRented}}
	_, err := l.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (l *LibraryMongoDb) GetAllbyAuthor(ctx context.Context, author string) ([]models.Book, error) {
	cursor, err := l.Collection.Find(ctx, bson.M{"author": author})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var books []models.Book
	for cursor.Next(ctx) {
		var book models.Book
		cursor.Decode(&book)
		books = append(books, book)
	}
	return books, nil
}

func (l *LibraryMongoDb) DeleteBookByISBN(ctx context.Context, isbn string) error {
	_, err := l.Collection.Find(ctx, bson.M{"isbn": isbn})
	if err != nil {
		return err
	}

	return nil
}

func (l *LibraryMongoDb) AggregateBooksByAuthor(ctx context.Context) ([]bson.M, error) {
	matchStage := bson.D{primitive.E{Key: "$match", Value: bson.D{}}}
	groupStage := bson.D{
		primitive.E{Key: "$group", Value: bson.D{
			primitive.E{Key: "_id", Value: "$author"},
			primitive.E{Key: "books", Value: bson.D{primitive.E{Key: "$push", Value: "$$ROOT"}}},
		}},
	}
	cursor, err := l.Collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []bson.M
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
