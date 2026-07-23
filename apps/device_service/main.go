package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
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

type CommandDefinition struct {
	Name             string         `json:"name" bson:"name"`
	ParametersSchema map[string]any `json:"parameters_schema" bson:"parameters_schema"`
}

type Device struct {
	ID       string              `json:"id" bson:"_id"`
	Name     string              `json:"name" bson:"name"`
	Type     string              `json:"type" bson:"type"`
	Settings map[string]any      `json:"settings" bson:"settings"`
	Commands []CommandDefinition `json:"commands" bson:"commands"`
}

type DeviceRepository interface {
	FindAll(context.Context) ([]Device, error)
	FindByID(context.Context, string) (Device, error)
	Save(context.Context, Device) (Device, error)
}

type MongoDeviceRepository struct{ collection *mongo.Collection }

func (r *MongoDeviceRepository) FindAll(ctx context.Context) ([]Device, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	devices := make([]Device, 0)
	if err := cursor.All(ctx, &devices); err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *MongoDeviceRepository) FindByID(ctx context.Context, id string) (Device, error) {
	var device Device
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&device)
	return device, err
}

func (r *MongoDeviceRepository) Save(ctx context.Context, device Device) (Device, error) {
	_, err := r.collection.InsertOne(ctx, device)
	return device, err
}

type TelemetryClient struct {
	baseURL string
	client  *http.Client
}

func (c *TelemetryClient) Send(ctx context.Context, deviceID string, data map[string]any) error {
	payload, err := json.Marshal(map[string]any{"device_id": deviceID, "data": data})
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/telemetry", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		return errors.New("history service rejected telemetry")
	}
	return nil
}

type DeviceService struct {
	repository DeviceRepository
	telemetry  *TelemetryClient
}

func (s *DeviceService) ListDevices(ctx context.Context) ([]Device, error) {
	return s.repository.FindAll(ctx)
}

func (s *DeviceService) GetDevice(ctx context.Context, id string) (Device, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *DeviceService) RegisterDevice(ctx context.Context, device Device) (Device, error) {
	for _, command := range device.Commands {
		if command.Name == "" || command.ParametersSchema == nil {
			return Device{}, errors.New("each command requires name and parameters_schema")
		}
	}
	device.ID = newID()
	return s.repository.Save(ctx, device)
}

func (s *DeviceService) ValidateCommand(ctx context.Context, id, command string) error {
	device, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	for _, supported := range device.Commands {
		if supported.Name == command {
			return nil
		}
	}
	return errors.New("command is not supported by device")
}

func (s *DeviceService) ExecuteCommand(ctx context.Context, id, command string, parameters map[string]any) error {
	if err := s.ValidateCommand(ctx, id, command); err != nil {
		return err
	}
	return s.telemetry.Send(ctx, id, map[string]any{
		"command": command, "parameters": parameters, "status": "executed",
	})
}

type DeviceController struct{ service *DeviceService }

func (h *DeviceController) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/devices", h.ListDevices)
	router.GET("/devices/:id", h.GetDevice)
	router.POST("/devices", h.RegisterDevice)
	router.POST("/devices/:id/commands", h.ExecuteCommand)
}

func (h *DeviceController) ListDevices(c *gin.Context) {
	devices, err := h.service.ListDevices(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (h *DeviceController) GetDevice(c *gin.Context) {
	device, err := h.service.GetDevice(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (h *DeviceController) RegisterDevice(c *gin.Context) {
	var device Device
	if err := c.ShouldBindJSON(&device); err != nil || device.Name == "" || device.Type == "" || device.Settings == nil || device.Commands == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name, type, settings and commands are required"})
		return
	}
	created, err := h.service.RegisterDevice(c.Request.Context(), device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Header("Location", "/api/devices/"+created.ID)
	c.JSON(http.StatusCreated, created)
}

func (h *DeviceController) ExecuteCommand(c *gin.Context) {
	var request struct {
		Command    string         `json:"command" binding:"required"`
		Parameters map[string]any `json:"parameters"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.ExecuteCommand(c.Request.Context(), c.Param("id"), request.Command, request.Parameters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
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

	repository := &MongoDeviceRepository{mongoClient.Database("devices").Collection("devices")}
	telemetry := &TelemetryClient{env("HISTORY_SERVICE_URL", "http://localhost:8083"), &http.Client{Timeout: 5 * time.Second}}
	controller := &DeviceController{&DeviceService{repository, telemetry}}
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	controller.RegisterRoutes(router.Group(""))
	log.Fatal(router.Run(":8082"))
}
