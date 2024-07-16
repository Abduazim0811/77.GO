package main

import (
	"Library/api"
	"Library/internal/mongodb"
	"log"
)


func main(){
	db, err := mongodb.MongoDb("mongodb://localhost:27017", "library", "books")
	if err != nil{
		log.Fatal(err)
	}
	api.Routes(db)
}