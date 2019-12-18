package model

import "go.mongodb.org/mongo-driver/bson/primitive"

//Car ...
type Car struct {
	IDMongo    primitive.ObjectID `bson:"_id" json:"id"`
	ID         string             `bson:"id" json:"ID"`
	Model      string             `bson:"model" json:"model"`
	Date       string             `bson:"date" json:"date"`
	// Membership struct {
	// 	GroupName string `bson:"groupname" json:"groupname"`
	// 	GroupID   int64  `bson:"groupid" json:"groupid"`
	// }
}

// {
// 	"id": "5df8c8974456199834dccc52",
// 	"ID": "5555 AB-5",
// 	"model": "Bently",
// 	"date": "01.01.1990",
// 	"Membership": {
// 		"groupname": "car",
// 		"groupid": 0
// 	}
// }
