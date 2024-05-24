/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/maguec/AcmeEverythingCompany/utils"
	"github.com/spf13/cobra"
)

var benchmarkType string
var benchmarkRuns int

// benchmarkCmd represents the benchmark command
var benchmarkCmd = &cobra.Command{
	Use:   "benchmark",
	Short: "Benchmark the database to show tuning effects",
	Long: `This runs some benchmarks on the database to show tuning effects

For options that are available in postgres and alloydb.`,
	Run: func(cmd *cobra.Command, args []string) {
		if benchmarkType == "catalog" || benchmarkType == "customer" {
			db, err := utils.GetDb(&Config)
			if err != nil {
				log.Fatal(err)
			}

			utils.Benchmark(db, benchmarkType, benchmarkRuns)
		} else {
			log.Fatal("Invalid benchmark type must be catalog or product")
		}
	},
}

func init() {
	rootCmd.AddCommand(benchmarkCmd)
	benchmarkCmd.Flags().StringVarP(&benchmarkType, "benchmark-type", "t", "catalog", "What type of benchmark to run catlog or customer")
	benchmarkCmd.Flags().IntVarP(&benchmarkRuns, "benchmark-runs", "r", 10, "How many times to run the benchmark")
}
