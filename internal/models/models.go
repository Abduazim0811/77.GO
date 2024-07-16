package models

type Book struct{
	ISBN 	string		`bson:"isbn" json:"isbn"`
	Title	string		`bson:"title" json:"title"`
	Author  string		`bson:"author" json:"author"`
	IsRented bool 		`bson:"is_rented" json:"is_rented"`
}