package utils

import (
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

type Customer struct {
	Id         uuid.UUID `gorm:"type:uuid"`
	FirstName  string    `fake:"{firstname}"`
	Lastname   string    `fake:"{lastname}"`
	Email      string    `fake:"{email}"`
	Company    string    `fake:"{company}"`
	JobTitle   string    `fake:"{jobtitle}"`
	Phone      string    `fake:"{phone}"`
  CreatedAt  time.Time
  UpdatedAt  time.Time
}

type Customers []Customer

func (c Customers) Generate(l int) Customers {
	f := gofakeit.New(0)
	var r Customers
	var p Customer
	for i := 0; i < l; i++ {
		err := f.Struct(&p)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, p)
	}
	return r
}

func (c Customers) Emails(db *gorm.DB) []string {
	var r []string
	db.Clauses(
		hints.CommentAfter("select", "controller='customers',action='emails',application='acme'"),
	).Model(&Customer{}).Pluck("Email", &r)
	return r
}

func (c Customers) IDs(db *gorm.DB, pluck bool) []uuid.UUID {
	var r []uuid.UUID
	var customers []Customer
	if pluck {
		db.Clauses(
			hints.CommentAfter("select", "controller='customers',action='IDs-pluck',application='acme'"),
		).Model(&Customer{}).Pluck("Id", &r)
	} else {
		db.Clauses(
			hints.CommentAfter("select", "controller='customers',action='IDs',application='acme'"),
		).Find(&customers)
		for _, c := range customers {
			r = append(r, c.Id)
		}
	}
	return r
}

func (c Customers) DbLoad(db *gorm.DB) error {
	bar := progressbar.NewOptions(len(c), progressbar.OptionSetDescription("Customer Loading"))
	var data []Customer
	var err error
	db.AutoMigrate(&Customer{})
	for i := 0; i < len(c); i++ {
		data = append(data, c[i])
		if len(data) == 100 {
			err = db.Clauses(
				hints.CommentAfter(
					"insert", "controller='catalog',action='DbLoad',application='acme'"),
			).Clauses(
				clause.OnConflict{UpdateAll: true},
			).Create(&data).Error
			if err != nil {
				return err
			}
			data = nil
			bar.Add(100)
		}
	}
	if len(data) > 0 {
		err = db.Clauses(
			hints.CommentAfter("insert", "controller='catalog',action='DbLoad',application='acme'"),
		).Clauses(
			clause.OnConflict{UpdateAll: true},
		).Create(&data).Error
		if err != nil {
			return err
		}
		bar.Add(len(data))
	}
	return err
}
