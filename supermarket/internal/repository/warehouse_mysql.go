package repository

import (
	"app/internal"
	"database/sql"
)

// RepositoryWarehouseMySQL is the repository warehouse MySQL.
type RepositoryWarehouseMySQL struct {
	// db is the connection to the database.
	db *sql.DB
}

// NewRepositoryWarehouseMySQL creates a new repository warehouse MySQL.
func NewRepositoryWarehouseMySQL(db *sql.DB) (rw *RepositoryWarehouseMySQL) {
	rw = &RepositoryWarehouseMySQL{
		db: db,
	}
	return
}

// FindById returns a warehouse by its id.
func (rw *RepositoryWarehouseMySQL) FindById(id int) (w internal.Warehouse, err error) {
	// Query the database for the warehouse.
	row := rw.db.QueryRow("SELECT id, name, address, telephone, capacity FROM warehouses WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return w, err
	}

	// Scan the row into the warehouse.
	err = row.Scan(&w.Id, &w.Name, &w.Address, &w.Telephone, &w.Capacity)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrRepositoryWarehouseNotFound
		}
		return
	}
	return
}

// GetAll returns all warehouses.
func (rw *RepositoryWarehouseMySQL) GetAll() (w []internal.Warehouse, err error) {
	// Query the database for the warehouses.
	rows, err := rw.db.Query("SELECT id, name, address, telephone, capacity FROM warehouses")
	if err != nil {
		return
	}
	defer rows.Close()

	// Scan the rows into the warehouses.
	for rows.Next() {
		var warehouse internal.Warehouse
		err = rows.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Address, &warehouse.Telephone, &warehouse.Capacity)
		if err != nil {
			return
		}
		w = append(w, warehouse)
	}
	return
}

// Create creates a warehouse.
func (rw *RepositoryWarehouseMySQL) Save(w *internal.Warehouse) (err error) {
	result, err := rw.db.Exec(
		"INSERT INTO warehouses (name, address, telephone, capacity) VALUES (?, ?, ?, ?)",
		w.Name, w.Address, w.Telephone, w.Capacity,
	)
	if err != nil {
		return
	}

	// Get the last inserted id.
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	// Set the id of the warehouse.
	w.Id = int(id)
	return
}
