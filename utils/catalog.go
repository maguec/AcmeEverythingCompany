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

type CatalogOrder struct {
	Id        uuid.UUID
	UnitPrice float64
}

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

func (c Catalog) IDs(db *gorm.DB) []CatalogOrder {
	var r []CatalogOrder
	var catalog Catalog
	db.Clauses(hints.CommentAfter("where", "type='products',func='IDs'")).Find(&catalog)
	for _, c := range catalog {
		r = append(r, CatalogOrder{Id: c.Id, UnitPrice: c.UnitPrice})
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
