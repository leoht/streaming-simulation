package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"leohetsch.com/simulation/producer"
	"leohetsch.com/simulation/simulation"
)

type JsonUserSimulation struct {
	UserId  string `json:"user_id"`
	Running bool   `json:"running"`
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Starting Kafka producer...")

	kafkaProducer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVER_URL"),

		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "OAUTHBEARER",
		"client.id":         "simulation-producer",
		"acks":              "all",
	})

	if err != nil {
		fmt.Printf("Failed to create producer: %s", err)
		os.Exit(1)
	}

	sim := simulation.StartSimulation()
	client := producer.NewKafkaClient(os.Getenv("TOPIC_NAME"), kafkaProducer)
	go client.Start(sim.ProducerChannel())

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Returns a list of current user simulations
	r.GET("/simulations", func(c *gin.Context) {
		simulations := simulation.AllUserSimulations()
		jsonSimulations := make([]JsonUserSimulation, len(simulations))
		for _, s := range simulations {
			jsonSimulations = append(jsonSimulations, JsonUserSimulation{s.UserId, s.Running})
		}

		c.JSON(http.StatusOK, gin.H{
			"simulations": simulations,
		})
	})

	// Start a new user simulation
	r.POST("/simulations", func(c *gin.Context) {
		simulation.StartNewUserSimulation()

		simulations := simulation.AllUserSimulations()
		// jsonSimulations := make([]JsonUserSimulation, len(simulations))
		// for _, s := range simulations {
		// 	jsonSimulations = append(jsonSimulations, JsonUserSimulation{s.UserId, s.Running})
		// }

		c.JSON(http.StatusOK, gin.H{
			"simulations": simulations,
		})
	})

	r.PUT("/simulations/:userId/stop", func(c *gin.Context) {
		userId := c.Param("userId")
		sim := simulation.StopSimulationForUser(userId)
		allSimulations := simulation.AllUserSimulations()
		if sim != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulations": allSimulations,
			})
		}
	})

	r.PUT("/simulations/:userId/resume", func(c *gin.Context) {
		userId := c.Param("userId")
		sim := simulation.ResumeSimulationForUser(userId)
		allSimulations := simulation.AllUserSimulations()
		if sim != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulations": allSimulations,
			})
		}
	})

	r.Run()
}
