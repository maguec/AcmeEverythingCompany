package utils

import (
	"fmt"
	"time"

	"github.com/jamiealquiza/tachymeter"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
)

func Benchmark(db *gorm.DB, benchmarktype string, loops int) error {
	var err error
	bar := progressbar.NewOptions(loops, progressbar.OptionSetDescription(fmt.Sprintf("Benchmarking %s", benchmarktype)))
	tach := tachymeter.New(&tachymeter.Config{Size: loops})
	catalog := Catalog{}
	customers := Customers{}
	for i := 0; i < loops; i++ {
		startTime := time.Now()
    if benchmarktype == "catalog" {
      catalog.IDs(db)
    } else if benchmarktype == "customer" {
      customers.IDs(db, true)
    }
  
		catalog.IDs(db)
		bar.Add(1)
		tach.AddTime(time.Since(startTime))
	}
	results := tach.Calc()
	fmt.Println("\n------------------ Latency ------------------")
	fmt.Printf(
		"Max:\t\t%s\nMin:\t\t%s\nP95:\t\t%s\nP99:\t\t%s\nP99.9:\t\t%s\n\n",
		results.Time.Max,
		results.Time.Min,
		results.Time.P95,
		results.Time.P99,
		results.Time.P999,
	)
	return err
}
