package main

import (
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/temperature", func(c *gin.Context) {
		sensorID := c.Query("sensorId")
		location := c.Query("location")

		if sensorID == "" {
			switch location {
			case "Living Room":
				sensorID = "1"
			case "Bedroom":
				sensorID = "2"
			case "Kitchen":
				sensorID = "3"
			default:
				sensorID = "0"
			}
		}

		if location == "" {
			switch sensorID {
			case "1":
				location = "Living Room"
			case "2":
				location = "Bedroom"
			case "3":
				location = "Kitchen"
			default:
				location = "Unknown"
			}
		}

		min, max := -10, 40
		value := min + rand.IntN(max-min+1)

		c.JSON(http.StatusOK, TemperatureResponse{
			Value:       float64(value),
			Unit:        "celsius",
			Timestamp:   time.Now(),
			Location:    location,
			Status:      "status",
			SensorID:    sensorID,
			SensorType:  "sensor-type",
			Description: "description",
		})
	})

	if err := router.Run(":8081"); err != nil {
		panic(err)
	}
}
