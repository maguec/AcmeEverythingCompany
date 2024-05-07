/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/maguec/AcmeEverythingCompany/utils"
	"github.com/spf13/cobra"
)

var maxClients, orderCount int

// addOrdersCmd represents the addOrders command
var addOrdersCmd = &cobra.Command{
	Use:   "addOrders",
	Short: "Insert a number of random orders",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		orders := utils.Orders{}
		orders.DbLoad(Config, orderCount, maxClients)

	},
}

func init() {
	rootCmd.AddCommand(addOrdersCmd)
	addOrdersCmd.Flags().IntVarP(&orderCount, "order-count", "o", Config.OrderCount, "How many orders to add to the database")
	addOrdersCmd.Flags().IntVarP(&maxClients, "max-clients", "m", Config.MaxClients, "How many simulataneous clients to run")

	// addOrdersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
