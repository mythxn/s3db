package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"s3db/internal"
)

type App struct {
	s3Config internal.S3Config
}

func main() {
	app := NewApp()

	r := gin.Default()
	r.GET("/records/:id", app.handleGetRecord)
	r.POST("/records/:id", app.handlePostRecord)
	r.GET("/records", app.handleGetRecords)
	r.POST("/drop-db", app.handleDropDB)

	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func NewApp() *App {
	return &App{
		s3Config: internal.S3Config{
			Region:          os.Getenv("S3_REGION"),
			AccessKeyID:     os.Getenv("S3_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("S3_SECRET_ACCESS_KEY"),
			BucketName:      os.Getenv("S3_BUCKET_NAME"),
		},
	}
}

func (app *App) handleGetRecord(c *gin.Context) {
	key := c.Param("id")
	resp, err := internal.GetRecord(app.s3Config, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var jsonData interface{}
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jsonData)
}

func (app *App) handlePostRecord(c *gin.Context) {
	key := c.Param("id")
	serializedBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData interface{}
	err = json.Unmarshal(serializedBody, &jsonData)
	if err != nil {
		c.String(http.StatusBadRequest, "got a non-JSON body")
		return
	}

	err = internal.NewRecord(app.s3Config, key, string(serializedBody))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusCreated)
}

func (app *App) handleGetRecords(c *gin.Context) {
	allObjectsList, err := internal.ListAllObjects(app.s3Config)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, allObjectsList)
}

func (app *App) handleDropDB(c *gin.Context) {
	err := internal.DropDB(app.s3Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Database dropped successfully"})
}
