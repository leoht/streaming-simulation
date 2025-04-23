package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "simulation",
	Short: "A streaming simulation",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("please provide a subcommand")
		}
		if args[0] == "generate-user-ids" {
			uuids := generateUserUuids(100)
			file, err := os.Create("users.txt")
			if err != nil {
				log.Fatal("could not create users.txt")
			}
			defer file.Close()

			text := uuidsToText(uuids)
			file.Write([]byte(text))

			fmt.Println("Saved list of generated user ids into users.txt")
		}
		if args[0] == "start" {
			// producer.Start()
		}
	},
}

// Generates a list of UUIDs of the given length
func generateUserUuids(length int) []uuid.UUID {
	list := make([]uuid.UUID, length)
	for i := 0; i < length; i++ {
		list[i] = uuid.New()
	}

	return list
}

func uuidsToText(list []uuid.UUID) string {
	var sb strings.Builder
	for _, uuid := range list {
		sb.WriteString(uuid.String())
		sb.WriteRune('\n')
	}
	return sb.String()
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
