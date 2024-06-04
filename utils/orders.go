package utils

import (
	"fmt"
	"log"
	"sync"
	"time"

	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

type Order struct {
	Id         uuid.UUID `gorm:"type:uuid"`
	OrderId    uuid.UUID `gorm:"type:uuid"`
	ProductId  uuid.UUID `gorm:"type:uuid"`
	CustomerId uuid.UUID `gorm:"type:uuid"`
	Units      int
	TotalCost  float64 `sql:"type:decimal(10,2);"`
	CreatedAt   time.Time 
	UpdatedAt   time.Time 
}

type Orders []Order

func writeOrders(
	w int,
	wg *sync.WaitGroup,
	cfg AcmeConfig,
	txns <-chan Order,
  bar *progressbar.ProgressBar,
) error {
	var err error
	var data []Order
	db, err := GetDb(&cfg)
	if err != nil {
		return err
	}
	defer wg.Done()
	for {
		select {
		case job, ok := <-txns:
			if !ok {
				fmt.Printf("Error with thread %d\n", w)
			}
			data = append(data, job)
			if len(data) == 100 {
				err = db.Clauses(
					hints.CommentAfter("insert", "controller='order',action='DbLoad',application='acme'"),
				).Clauses(
					clause.OnConflict{UpdateAll: true},
				).Create(&data).Error
				if err != nil {
					return err
				}
				data = nil
        bar.Add(100)
			}

		default:
			if len(data) > 0 {
				err = db.Clauses(
					hints.CommentAfter("insert", "controller='order',action='DbLoad',application='acme'"),
				).Clauses(
					clause.OnConflict{UpdateAll: true},
				).Create(&data).Error
				if err != nil {
					return err
				}
				data = nil
        bar.Add(len(data))
			}
			return err
		}
	}
}

func (o Orders) DbLoad(cfg AcmeConfig, orderCount, maxClients int) error {
	var err error
	db, err := GetDb(&cfg)
  // Create the table if it doesn't exist
  err = db.Migrator().AutoMigrate(&Order{})
	if err != nil {
		log.Fatal(err)
	}
	customers := Customers{}
	customerUUIDS := customers.IDs(db, false)
	catalog := Catalog{}
	productUUIDS := catalog.IDs(db)
	orders := []Order{}
  orderbar := progressbar.NewOptions(orderCount, progressbar.OptionSetDescription("Generating Orders"))
	for i := 0; i < orderCount; i++ {
		myuuid := uuid.New()
		customerId := customerUUIDS[rand.Intn(len(customerUUIDS))]
		for j := 0; j < rand.Intn(15)+1; j++ {
			z := productUUIDS[rand.Intn(len(productUUIDS))]
			units := rand.Intn(50) + 1
			o_uuid := uuid.New()
      t := gofakeit.DateRange(time.Now().AddDate(0, 0, -180), time.Now())
			ord := Order{
				OrderId:    myuuid,
				Id:         o_uuid,
				CustomerId: customerId,
				ProductId:  z.Id,
				Units:      units,
				TotalCost:  z.UnitPrice * float64(units),
        CreatedAt:  t,
        UpdatedAt:  t,
			}
			orders = append(orders, ord)
		}
    orderbar.Add(1)
	}
  bar := progressbar.NewOptions(len(orders), progressbar.OptionSetDescription("Orders Loading"))
	txns := make(chan Order, len(orders))
	for _, o := range orders {
		txns <- o
	}
	wg := new(sync.WaitGroup)
	for w := 0; w < maxClients; w++ {
		wg.Add(1)
		go writeOrders(w, wg, cfg, txns, bar)
	}
	wg.Wait()
  fmt.Println("") // blank line for readability
	return err
}
