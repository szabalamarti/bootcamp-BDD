package internal

import "errors"

// Warehouse is a struct that contains the attributes of a warehouse
type Warehouse struct {
	// Id is the unique identifier of the warehouse
	Id int
	// Name is the name of the warehouse
	Name string
	// Address is the address of the warehouse
	Address string
	// Telephone is the telephone of the warehouse
	Telephone string
	// Capacity is the capacity of the warehouses
	Capacity int
}

// WarehouseReportProducts is a struct that contains the amount of products in a warehouse
type WarehouseReportProducts struct {
	// WarehouseName is the name of the warehouse
	WarehouseName string
	// ProductCount is the amount of products in the warehouse
	ProductCount int
}

var (
	// ErrRepositoryWarehouseNotFound is returned when a warehouse is not found.
	ErrRepositoryWarehouseNotFound = errors.New("repository: warehouse not found")
)

// RepositoryWarehouse is an interface that contains the methods for a warehouse repository
type RepositoryWarehouse interface {
	// FindById returns a warehouse by its id
	FindById(id int) (w Warehouse, err error)
	// GetAll returns all warehouses
	GetAll() (w []Warehouse, err error)
	// Save creates a warehouse
	Save(w *Warehouse) (err error)
	// ReportProducts returns a report of the amount of products in each warehouse
	ReportProducts(warehouseIds []int) (r []WarehouseReportProducts, err error)
}
