package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"leohetsch.com/simulation/producer"
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

	go producer.Start()

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Returns a list of current user simulations
	r.GET("/simulations", func(c *gin.Context) {
		simulations := producer.GetSimulations()
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
		producer.StartNewSimulation()

		simulations := producer.GetSimulations()
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
		simulation := producer.StopSimulationForUser(userId)
		simulations := producer.GetSimulations()
		if simulation != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulations": simulations,
			})
		}
	})

	r.PUT("/simulations/:userId/resume", func(c *gin.Context) {
		userId := c.Param("userId")
		simulation := producer.ResumeSimulationForUser(userId)
		simulations := producer.GetSimulations()
		if simulation != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulations": simulations,
			})
		}
	})

	r.Run()
}
