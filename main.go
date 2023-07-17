package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"s3db/internal"
)

func main() {
	s3Config := internal.S3Config{
		Region:          os.Getenv("S3_REGION"),
		AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
		BucketName:      os.Getenv("S3_BUCKET_NAME"),
	}

	r := gin.Default()

	r.POST("/drop-db", func(c *gin.Context) {
		err := internal.DropDB(s3Config)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Database dropped successfully"})
	})

	r.POST("/new-record", func(c *gin.Context) {
		var newRecordRequest internal.NewRecordRequest
		err := c.ShouldBindJSON(&newRecordRequest)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		err = internal.NewRecord(s3Config, newRecordRequest)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Record added successfully"})
	})

	r.GET("/get-record", func(c *gin.Context) {
		var getRecordRequest internal.GetRecordRequest
		err := c.ShouldBindJSON(&getRecordRequest)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		value, err := internal.GetRecord(s3Config, getRecordRequest)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"value": value})
	})

	err := r.Run()
	if err != nil {
		panic(err)
	}
}
