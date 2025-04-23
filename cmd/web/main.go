package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"leohetsch.com/simulation/producer"
)

type JsonUserSimulation struct {
	UserId  string `json:"user_id"`
	Running bool   `json:"running"`
}

func main() {
	go producer.Start()

	r := gin.Default()
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

		fmt.Println(simulations)
		c.JSON(http.StatusOK, gin.H{
			"simulations": simulations,
		})
	})

	r.PUT("/simulations/:userId/stop", func(c *gin.Context) {
		userId := c.Param("userId")
		simulation := producer.StopSimulationForUser(userId)
		if simulation != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulation": JsonUserSimulation{simulation.UserId, simulation.Running},
			})
		}
	})

	r.PUT("/simulations/:userId/resume", func(c *gin.Context) {
		userId := c.Param("userId")
		simulation := producer.ResumeSimulationForUser(userId)
		if simulation != nil {
			c.JSON(http.StatusOK, gin.H{
				"simulation": JsonUserSimulation{simulation.UserId, simulation.Running},
			})
		}
	})

	r.Run()
}
