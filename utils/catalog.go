package utils

import (
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/hints"
)

type Product struct {
	gorm.Model  `fake:"skip"`
	Id          uuid.UUID `gorm:"type:uuid"`
	UnitPrice   float64   `sql:"type:decimal(10,2);" fake:"{price}"`
	Name        string    `fake:"{productname}"`
	Description string    `fake:"{productdescription}"`
	Category    string    `fake:"{productcategory}"`
}

type Catalog []Product

func (c Catalog) Generate(l int) Catalog {
	f := gofakeit.New(0)
	var r Catalog
	var p Product
	for i := 0; i < l; i++ {
		err := f.Struct(&p)
		if err != nil {
			log.Fatal(err)
		}
		r = append(r, p)
	}
	return r
}

func (c Catalog) IDs() []uuid.UUID {
	var r []uuid.UUID
	for i := 0; i < len(c); i++ {
		r = append(r, c[i].Id)
	}
	return r
}

func (c Catalog) DbLoad(db *gorm.DB) error {
	var data []Product
	var err error
	db.AutoMigrate(&Product{})
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
