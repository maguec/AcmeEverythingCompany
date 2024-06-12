/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
  "github.com/maguec/AcmeEverythingCompany/utils"
)

// analyticsCmd represents the analytics command
var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Get timing information for analytical queries",
	Long: `Runs several queries to get timing information for analytical queries
This is generally run to show tuning improvments`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := utils.GetDb(&Config)
    if err != nil {
      log.Fatal(err)
    }
    utils.Analytics(db)
	},
}

func init() {
	rootCmd.AddCommand(analyticsCmd)
}
