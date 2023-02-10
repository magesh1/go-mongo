package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Baseurl = "localhost:8080/"

// request
type ShortenBody struct {
	OriginalUrl string `json:"original_url"`
}

// struct to store in db
type UrlDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Urlcode     string             `bson:"urlcode"`
	Originalurl string             `bson:"originalurl"`
	ShortUrl    string             `bson:"shorturl"`
	CreatedAt   time.Time          `bson:"createdAt"`
	ExpiresAt   time.Time          `bson:"expiresAt"`
}
