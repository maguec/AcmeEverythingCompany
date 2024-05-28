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

    if customerCount < 1 {
      customerCount = Config.CustomerCount
    }
    if productCount < 1 {
      productCount = Config.ProductCount
    }
		if debug {
			fmt.Printf("Creating %d customers and %d products\n", customerCount, productCount)
		}
		catalog := utils.Catalog{}
		catalog = catalog.Generate(productCount)
		customers := utils.Customers{}
		customers = customers.Generate(customerCount)
		db, err := utils.GetDb(&Config)
		if err != nil {
			log.Fatal(err)
		}
		catalog.DbLoad(db)
		customers.DbLoad(db)
		fmt.Println("") // Blank line to that it looks nice
		if debug {
			fmt.Printf("\nDatabase seeded with %d customers and %d products\n", customerCount, productCount)
		}

	},
}

func init() {
	rootCmd.AddCommand(dbseedCmd)
  dbseedCmd.Flags().IntVarP(&customerCount, "customer-count", "c", 0, fmt.Sprintf("Override the number of customers to add to the database default: %d", Config.CustomerCount))
  dbseedCmd.Flags().IntVarP(&productCount, "product-count", "p", 0, fmt.Sprintf("Override the number of products to add to the database default: %d",  Config.ProductCount))
}
