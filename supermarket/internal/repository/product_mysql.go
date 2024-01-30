package repository

import (
	"database/sql"
	"errors"
	"supermarket/internal"

	"github.com/go-sql-driver/mysql"
)

// NewRepositoryProductMySQL creates a new repository product MySQL.
func NewRepositoryProductMySQL(db *sql.DB) (rp *RepositoryProductMySQL) {
	rp = &RepositoryProductMySQL{
		db: db,
	}
	return
}

// RepositoryProductMySQL is the repository product MySQL.
type RepositoryProductMySQL struct {
	// db is the connection to the database.
	db *sql.DB
}

// FindById returns a product by its id.
func (rp *RepositoryProductMySQL) FindById(id int) (p internal.Product, err error) {
	// Query the database for the product.
	row := rp.db.QueryRow("SELECT id, name, price, quantity, code_value, is_published, expiration, price, id_warehouse FROM products WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return p, err
	}

	// Scan the row into the product.
	err = row.Scan(&p.Id, &p.Name, &p.Price, &p.Quantity, &p.CodeValue, &p.IsPublished, &p.Expiration, &p.Price, &p.WarehouseId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = internal.ErrRepositoryProductNotFound
		}
		return
	}
	return
}

// GetAll returns all products.
func (rp *RepositoryProductMySQL) GetAll() (p []internal.Product, err error) {
	// Query the database for the products.
	rows, err := rp.db.Query("SELECT id, name, price, quantity, code_value, is_published, expiration, price, id_warehouse FROM products")
	if err != nil {
		return
	}
	defer rows.Close()

	// Scan the rows into the products.
	for rows.Next() {
		var product internal.Product
		err = rows.Scan(&product.Id, &product.Name, &product.Price, &product.Quantity, &product.CodeValue, &product.IsPublished, &product.Expiration, &product.Price, &product.WarehouseId)
		if err != nil {
			return
		}
		p = append(p, product)
	}
	return
}

func (rp *RepositoryProductMySQL) Save(p *internal.Product) (err error) {
	result, err := rp.db.Exec(
		"INSERT INTO products (name, quantity, code_value, is_published, expiration, price, id_warehouse) VALUES (?, ?, ?, ?, ?, ?, ?)",
		p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.WarehouseId,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrBuyerRepositoryDuplicated
			default:
				// ...
			}
			return
		}

		return
	}

	// Get the last inserted id.
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	// Set the id of the product.
	p.Id = int(id)
	return
}

func (rp *RepositoryProductMySQL) UpdateOrSave(p *internal.Product) (err error) {
	err = rp.Update(p)
	if err == internal.ErrRepositoryProductNotFound {
		err = rp.Save(p)
	}
	return
}

func (rp *RepositoryProductMySQL) Update(p *internal.Product) (err error) {
	_, err = rp.db.Exec(
		"UPDATE products SET name = ?, quantity = ?, code_value = ?, is_published = ?, expiration = ?, price = ?, id_warehouse = ? WHERE id = ?",
		p.Name, p.Quantity, p.CodeValue, p.IsPublished, p.Expiration, p.Price, p.WarehouseId, p.Id,
	)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			switch mysqlErr.Number {
			case 1062:
				err = internal.ErrBuyerRepositoryDuplicated
			default:
				// ...
			}
			return
		}
	}

	return
}

func (rp *RepositoryProductMySQL) Delete(id int) (err error) {
	result, err := rp.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return
	}

	// check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return
	}
	// if no rows were affected, return not found error
	if rowsAffected == 0 {
		err = internal.ErrRepositoryProductNotFound
		return
	}

	return
}
