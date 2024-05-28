package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jamiealquiza/tachymeter"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
	"gorm.io/hints"
)

func hasIndex(db *gorm.DB, name string) bool {

  indexes, err := db.Migrator().GetIndexes(&Customer{})
	if err != nil {
		panic(err)
	}
	for _, index := range indexes {
    if index.Name() == name {
      return true
    }
	}
	return false
}

func indexWorker(
	t int,
	wg *sync.WaitGroup,
	jobs <-chan string,
	db *gorm.DB,
	bar *progressbar.ProgressBar,
	tach *tachymeter.Tachymeter,
  insight string,
) {
	defer wg.Done()
	for {
		select {
		case job, ok := <-jobs:
			if !ok {
				return
			}
			startTime := time.Now()
			e := db.Clauses(
				hints.CommentAfter("select", insight),
			).Model(&Customer{}).Where("Email = ?", job).Select("Email, Id").Find(&Customer{})
			if e.Error != nil {
				log.Fatal(t, e.Error)
			}
			bar.Add(1)
			tach.AddTime(time.Since(startTime))
		default:
			return
		}
	}
}

func IndexBenchmark(db *gorm.DB, loops int) error {
	var err error

  /**********************************************************************************************
    CONGRATULATIONS! 

    You've found the slow query add the folowing index to see the performance improvement:
    
    CREATE INDEX faster ON customers(email) INCLUDE(id);

    You will see the action TAG change after you re-run the benchmark

  **********************************************************************************************/

  insight := "controller='customers',action='benchmark-no-index',application='acme'"
  if hasIndex(db, "faster") {
    insight = "controller='customers',action='benchmark-with-index',application='acme'"
  }

	customers := Customers{}
	emails := customers.Emails(db)
	wg := sync.WaitGroup{}
	bar := progressbar.NewOptions(loops*len(emails), progressbar.OptionSetDescription("Benchmarking Index"))
	tach := tachymeter.New(&tachymeter.Config{Size: loops * len(emails)})
	jobs := make(chan string, loops*len(emails))
	for i := 0; i < loops; i++ {
		for j := 0; j < len(emails); j++ {
			jobs <- emails[j]
		}
	}
	for w := 0; w < 10; w++ {
		wg.Add(1)
		go indexWorker(w, &wg, jobs, db, bar, tach, insight)
	}
	wg.Wait()
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
