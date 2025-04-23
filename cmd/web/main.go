package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

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

func readUserIds() []string {
	contents, err := os.ReadFile("users.txt")
	if err != nil {
		log.Fatal("could not read users.txt ")
	}

	return splitLines(string(contents))
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

func main() {
	err := godotenv.Load()
	userIds := readUserIds()
	producerInChannel := make(chan simulation.Event)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	simulation.StartSimulation()
	go producer.Start(producerInChannel)

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Returns a list of current user simulations
	r.GET("/simulations", func(c *gin.Context) {
		simulations := simulation.GetSimulations()
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
		simulation.StartNewSimulation(userIds, producerInChannel)

		simulations := simulation.GetSimulations()
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
		allSimulations := simulation.GetSimulations()
		if sim != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulations": allSimulations,
			})
		}
	})

	r.PUT("/simulations/:userId/resume", func(c *gin.Context) {
		userId := c.Param("userId")
		sim := simulation.ResumeSimulationForUser(userId)
		allSimulations := simulation.GetSimulations()
		if sim != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulations": allSimulations,
			})
		}
	})

	r.Run()
}
