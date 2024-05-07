package utils

import (
	"fmt"
	"log"
	"sync"

	"math/rand"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

type Order struct {
	gorm.Model `fake:"skip"`
	Id         uuid.UUID `gorm:"type:uuid"`
	OrderId    uuid.UUID `gorm:"type:uuid"`
	ProductId  uuid.UUID `gorm:"type:uuid"`
	CustomerId uuid.UUID `gorm:"type:uuid"`
	Units      int
	TotalCost  float64 `sql:"type:decimal(10,2);"`
}

type Orders []Order

func writeOrders(
	w int,
	wg *sync.WaitGroup,
	cfg AcmeConfig,
	txns <-chan Order,
) error {
	var err error
	var data []Order
	db, err := GetDb(&cfg)
	if err != nil {
		return err
	}
	db.AutoMigrate(&Order{})
	defer wg.Done()
	for {
		select {
		case job, ok := <-txns:
			if !ok {
				fmt.Printf("Error with thread %d\n", w)
			}
			data = append(data, job)
			if len(data) == 2 {
				err = db.Clauses(hints.CommentAfter("returning", "type='order',func='DbLoad'")).Clauses(clause.OnConflict{UpdateAll: true}).Create(&data).Error
				if err != nil {
					return err
				}
				data = nil
			}

		default:
			if len(data) > 0 {
				err = db.Clauses(hints.CommentAfter("returning", "type='order',func='DbLoad'")).Clauses(clause.OnConflict{UpdateAll: true}).Create(&data).Error
				if err != nil {
					return err
				}
				data = nil
			}
			return err
		}
	}
}

func (o Orders) DbLoad(cfg AcmeConfig, orderCount, maxClients int) error {
	var err error
	db, err := GetDb(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	customers := Customers{}
	customerUUIDS := customers.IDs(db)
	fmt.Printf("Length of customers: %d\n", len(customerUUIDS))
	catalog := Catalog{}
	productUUIDS := catalog.IDs(db)
	orders := []Order{}
	for i := 0; i < orderCount; i++ {
		myuuid := uuid.New()
		customerId := customerUUIDS[rand.Intn(len(customerUUIDS))]
		for j := 0; j < rand.Intn(15)+1; j++ {
			z := productUUIDS[rand.Intn(len(customerUUIDS))]
			units := rand.Intn(50) + 1
			o_uuid := uuid.New()
			ord := Order{
				OrderId:    myuuid,
				Id:         o_uuid,
				CustomerId: customerId,
				ProductId:  z.Id,
				Units:      units,
				TotalCost:  z.UnitPrice * float64(units),
			}
			orders = append(orders, ord)
		}
	}
	txns := make(chan Order, len(orders))
	for _, o := range orders {
		txns <- o
	}
	wg := new(sync.WaitGroup)
	for w := 0; w < maxClients; w++ {
		wg.Add(1)
		go writeOrders(w, wg, cfg, txns)
	}
	wg.Wait()
	return err
}
