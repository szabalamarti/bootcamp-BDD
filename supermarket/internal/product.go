package internal

import (
	"errors"
	"time"
)

// ProductAttributes is a struct that contains the attributes of a product
type ProductAttributes struct {
	// Name is the name of the product
	Name string
	// Quantity is the quantity of the product
	Quantity int
	// CodeValue is the code value of the product
	CodeValue string
	// IsPublished is the published status of the product
	IsPublished bool
	// Expiration
	Expiration time.Time
	// Price
	Price float64
	// WarehouseId is the id of the warehouse where the product is stored
	WarehouseId int
}

// Product is a struct that contains the attributes of a product
type Product struct {
	// Id is the unique identifier of the product
	Id int
	// ProductAttributes is the attributes of the product
	ProductAttributes
}

var (
	// ErrRepositoryProductNotFound is returned when a product is not found.
	ErrRepositoryProductNotFound = errors.New("repository: product not found")
	// ErrBuyerRepositoryDuplicated is returned when a product is duplicated.
	ErrBuyerRepositoryDuplicated = errors.New("repository: product duplicated")
)

// RepositoryProduct is an interface that contains the methods for a product repository
type RepositoryProduct interface {
	// FindById returns a product by its id
	FindById(id int) (p Product, err error)
	// GetAll returns all products
	GetAll() (p []Product, err error)
	// Save saves a product
	Save(p *Product) (err error)
	// UpdateOrSave updates or saves a product
	UpdateOrSave(p *Product) (err error)
	// Update updates a product
	Update(p *Product) (err error)
	// Delete deletes a product
	Delete(id int) (err error)
}

// StoreProduct is an interface for a product store.
type StoreProduct interface {
	// ReadAll reads all products from the store.
	ReadAll() (p map[int]Product, err error)
	// WriteAll writes all products to the store.
	WriteAll(p map[int]Product) (err error)
}
