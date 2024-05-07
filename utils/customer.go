package utils

import (
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

type Customer struct {
	gorm.Model `fake:"skip"`
	Id         uuid.UUID `gorm:"type:uuid"`
	FirstName  string    `fake:"{firstname}"`
	Lastname   string    `fake:"{lastname}"`
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

func (c Customers) IDs() []uuid.UUID {
	var r []uuid.UUID
	for i := 0; i < len(c); i++ {
		r = append(r, c[i].Id)
	}
	return r
}

func (c Customers) DbLoad(db *gorm.DB) error {
	var data []Customer
	var err error
	db.AutoMigrate(&Customer{})
	for i := 0; i < len(c); i++ {
		data = append(data, c[i])
		if len(data) == 100 {
			err = db.Clauses(hints.CommentAfter("returning", "type='catalog',func='DbLoad'")).Clauses(clause.OnConflict{UpdateAll: true}).Create(&data).Error
			if err != nil {
				return err
			}
			data = nil
		}
	}
	if len(data) > 0 {
		err = db.Clauses(hints.CommentAfter("returning", "type='catalog',func='DbLoad'")).Clauses(clause.OnConflict{UpdateAll: true}).Create(&data).Error
		if err != nil {
			return err
		}
	}
	return err
}
