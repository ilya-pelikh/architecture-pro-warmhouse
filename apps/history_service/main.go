package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Telemetry struct {
	ID         string         `json:"id" bson:"_id"`
	DeviceID   string         `json:"device_id" bson:"device_id"`
	RecordedAt time.Time      `json:"recorded_at" bson:"recorded_at"`
	Data       map[string]any `json:"data" bson:"data"`
}

type TelemetryRepository struct{ collection *mongo.Collection }

func (r *TelemetryRepository) Save(ctx context.Context, item Telemetry) (Telemetry, error) {
	_, err := r.collection.InsertOne(ctx, item)
	return item, err
}

func (r *TelemetryRepository) FindAll(ctx context.Context, deviceID string) ([]Telemetry, error) {
	filter := bson.M{}
	if deviceID != "" {
		filter["device_id"] = deviceID
	}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "recorded_at", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	items := make([]Telemetry, 0)
	if err := cursor.All(ctx, &items); err != nil {
		return nil, err
	}
	return items, nil
}

type HistoryController struct{ repository *TelemetryRepository }

func (h *HistoryController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/telemetry", h.ListTelemetry)
	router.POST("/telemetry", h.CreateTelemetry)
}

func (h *HistoryController) ListTelemetry(c *gin.Context) {
	items, err := h.repository.FindAll(c.Request.Context(), c.Query("device_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *HistoryController) CreateTelemetry(c *gin.Context) {
	var item Telemetry
	if err := c.ShouldBindJSON(&item); err != nil || item.DeviceID == "" || item.Data == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device_id and data are required"})
		return
	}
	item.ID = newID()
	item.RecordedAt = time.Now().UTC()
	created, err := h.repository.Save(c.Request.Context(), item)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

func newID() string {
	value := make([]byte, 16)
	_, _ = rand.Read(value)
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	encoded := hex.EncodeToString(value)
	return fmt.Sprintf("%s-%s-%s-%s-%s", encoded[0:8], encoded[8:12], encoded[12:16], encoded[16:20], encoded[20:32])
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(env("MONGO_URL", "mongodb://localhost:27017")))
	if err != nil {
		log.Fatal(err)
	}
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	repository := &TelemetryRepository{mongoClient.Database("history").Collection("telemetries")}
	controller := &HistoryController{repository}
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	controller.RegisterRoutes(router.Group(""))
	log.Fatal(router.Run(":8083"))
}
