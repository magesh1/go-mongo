package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/magesh1/go-mongo/db"
	"github.com/magesh1/go-mongo/models"
	"github.com/teris-io/shortid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func ShortenUrl(c *gin.Context) {

	// select {
	// case middleware.Limit <- struct{}{}:
	// 	defer func() {
	// 		<-middleware.Limit
	// 	}()
	// default:
	// 	c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
	// 	return
	// }

	var body models.ShortenBody

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check https or http present in the url
	_, urlErr := url.ParseRequestURI(body.OriginalUrl)

	if urlErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": urlErr.Error(),
		})
		return
	}

	urlCode, idErr := shortid.Generate()

	if idErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": idErr.Error(),
		})
		return
	}

	var result bson.M

	queryErr := db.Collection.FindOne(context.Background(), bson.D{{Key: "urlCode", Value: urlCode}}).Decode(&result)

	// check if data already present in the db
	if queryErr != nil {
		// check whether data exists
		if queryErr != mongo.ErrNoDocuments {
			c.JSON(http.StatusInternalServerError, gin.H{"err: ": queryErr.Error()})
			return
		}

	}

	if len(result) > 0 {

		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint("url already exist: ", urlCode)})
		return
	}

	date := time.Now()
	expires := date.AddDate(0, 0, 2)
	newUrl := models.Baseurl + urlCode
	docId := primitive.NewObjectID()

	newDoc := &models.UrlDoc{
		ID:          docId,
		Urlcode:     urlCode,
		Originalurl: body.OriginalUrl,
		ShortUrl:    newUrl,
		CreatedAt:   time.Now(),
		ExpiresAt:   expires,
	}

	_, err := db.Collection.InsertOne(context.Background(), newDoc)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"newURL":  newUrl,
		"expires": expires.Format("2006-01-02 15:04:05"),
	})
}
