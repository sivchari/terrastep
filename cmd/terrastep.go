package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "terrastep",
}

// Run runs command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("run", "r", "", "run terrastep")
}
