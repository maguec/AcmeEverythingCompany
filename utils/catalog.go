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

type Product struct {
	Id          uuid.UUID `gorm:"type:uuid"`
	UnitPrice   float64   `sql:"type:decimal(10,2);" fake:"{price}"`
	Name        string    `fake:"{productname}"`
	Description string    `fake:"{productdescription}"`
	Category    string    `fake:"{productcategory}"`
	CreatedAt   time.Time `fake:"{daterange:2023-01-01,2023-12-31,2006-01-02}" format:"yyyy-MM-dd"`
	UpdatedAt   time.Time `fake:"{daterange:2024-01-01,2024-12-31,2006-01-02}" format:"yyyy-MM-dd"`
}

// This is a slim version of the product to improve performance
type SlimProduct struct {
	Id        uuid.UUID
	UnitPrice float64
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
	var slimproduct []SlimProduct
	db.Clauses(
		hints.CommentAfter("select", "controller='catalog',action='Slim-IDs',application='acme'"),
	).Model(&Product{}).Find(&slimproduct)
	for _, c := range slimproduct {
		r = append(r, CatalogOrder{Id: c.Id, UnitPrice: c.UnitPrice})
	}
	return r
}

func (c Catalog) DbLoad(db *gorm.DB) error {
	var data []Product
	var err error
	bar := progressbar.NewOptions(len(c), progressbar.OptionSetDescription("Catalog Loading"))
	db.AutoMigrate(&Product{})
	for i := 0; i < len(c); i++ {
		data = append(data, c[i])
		if len(data) == 100 {
			err = db.Clauses(
				hints.CommentAfter("insert", "controller='catalog',action='DbLoad',application='acme'"),
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
