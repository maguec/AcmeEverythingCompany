/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/maguec/AcmeEverythingCompany/utils"
	"github.com/spf13/cobra"
)

var productCount, customerCount int

// dbseedCmd represents the dbseed command
var dbseedCmd = &cobra.Command{
	Use:   "dbseed",
	Short: "Seed the databse with sample data",
	Long: `Add some initial data to the database
Including users, organizations, and a product catalog.`,
	Run: func(cmd *cobra.Command, args []string) {
		catalog := utils.Catalog{}
		catalog = catalog.Generate(productCount)
		customers := utils.Customers{}
		customers = customers.Generate(customerCount)
		db, err := utils.GetDb(&Config)
		if err != nil {
			log.Fatal(err)
		}
		if debug {
			fmt.Printf("Creating and loading %d customers and %d products\n", Config.CustomerCount, Config.ProductCount)
		}
		catalog.DbLoad(db)
		customers.DbLoad(db)
    fmt.Println("")
    if debug {
    fmt.Printf("\nDatabase seeded with %d customers and %d products\n", customerCount, productCount)
    }

	},
}

func init() {
	rootCmd.AddCommand(dbseedCmd)
	dbseedCmd.Flags().IntVarP(&customerCount, "customer-count", "c", Config.CustomerCount, "How many customers to add to the database")
	dbseedCmd.Flags().IntVarP(&productCount, "product-count", "p", Config.ProductCount, "How many products to add to the database")
}
