package producer

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
)

// Start random user simulations and record produced events into
// the PostgresSQL database (TODO)
func Start() {

	contents, err := os.ReadFile("users.txt")
	if err != nil {
		log.Fatal("could not read users.txt ")
	}

	userIds := splitLines(string(contents))

	// For now let's start just one user simulation.
	userId := userIds[rand.Intn(len(userIds))]
	simulation := NewUserSimulation(userId)
	simulation.Start(userId, []string{"sign_in", "sign_up"})
}

func splitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}
