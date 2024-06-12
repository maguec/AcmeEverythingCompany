package utils

import (
	"time"
  "fmt"
  "sort"

	"gorm.io/gorm"
  "github.com/schollz/progressbar/v3"
)


func Analytics(db *gorm.DB) {
  queries := make(map[string]string)
  timings := make(map[string]time.Duration)
  queries["Min_cost"] = "SELECT min(total_cost) FROM orders"
  queries["Max_cost"] = "SELECT max(total_cost) FROM orders"
  queries["Sum_cost"] = "SELECT max(total_cost) FROM orders"
  queries["Top10_customers"] = "SELECT customer_id, SUM(total_cost) AS total_spent FROM orders GROUP BY customer_id ORDER BY total_spent DESC LIMIT 10"
  queries["Top10_products"] = "SELECT product_id, SUM(total_cost) AS total_revenue FROM orders GROUP BY product_id ORDER BY total_revenue ASC LIMIT 10"
  queries["Average_order_size"] = "SELECT AVG(order_total) AS avg_order_total FROM ( SELECT order_id, SUM(total_cost) AS order_total FROM orders GROUP BY order_id) AS order_totals"
  queries["Average_units_per_order"] = "SELECT AVG(order_total) AS avg_order_total FROM ( SELECT order_id, SUM(total_cost) AS order_total FROM orders GROUP BY order_id) AS order_totals"

bar := progressbar.NewOptions(len(queries), progressbar.OptionSetDescription("Analytics Running"))

  for k, v := range queries {
    bar.Add(1)
    start := time.Now()
    db.Exec(v)
    timings[k] = time.Since(start)
  }

  keys := []string{}
  for k := range timings {
    keys = append(keys, k)
  }

  sort.Strings(keys)

  fmt.Println("") //empty line
  fmt.Println("-------------------------------------")

  total := time.Duration(0)
  for _, k := range keys {
    v := timings[k]
    fmt.Printf("%-25s: %-5d (ms)\n", k, v.Milliseconds())
    total += v
  }
  fmt.Println("-------------------------------------")
  fmt.Printf("%-25s: %-5d (ms)\n", "TOTAL", total.Milliseconds())
}
