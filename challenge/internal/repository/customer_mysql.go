package repository

import (
	"database/sql"
	"math"

	"app/internal"
)

// NewCustomersMySQL creates new mysql repository for customer entity.
func NewCustomersMySQL(db *sql.DB) *CustomersMySQL {
	return &CustomersMySQL{db}
}

// CustomersMySQL is the MySQL repository implementation for customer entity.
type CustomersMySQL struct {
	// db is the database connection.
	db *sql.DB
}

// FindAll returns all customers from the database.
func (r *CustomersMySQL) FindAll() (c []internal.Customer, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `id`, `first_name`, `last_name`, `condition` FROM customers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var cs internal.Customer
		// scan the row into the customer
		err := rows.Scan(&cs.Id, &cs.FirstName, &cs.LastName, &cs.Condition)
		if err != nil {
			return nil, err
		}
		// append the customer to the slice
		c = append(c, cs)
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// Save saves the customer into the database.
func (r *CustomersMySQL) Save(c *internal.Customer) (err error) {
	// execute the query
	res, err := r.db.Exec(
		"INSERT INTO customers (`first_name`, `last_name`, `condition`) VALUES (?, ?, ?)",
		(*c).FirstName, (*c).LastName, (*c).Condition,
	)
	if err != nil {
		return err
	}

	// get the last inserted id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set the id
	(*c).Id = int(id)

	return
}

// FindTotalByCondition returns the aggregated money from invoices by customer condition.
// values rounded to the second decimal place.
func (r *CustomersMySQL) FindTotalByCondition() (t []internal.TotalByCondition, err error) {
	// execute the query
	rows, err := r.db.Query("SELECT `condition`, SUM(`total`) FROM customers INNER JOIN invoices ON customers.id = invoices.customer_id GROUP BY `condition`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var tb internal.TotalByCondition
		// scan the row into the total by condition
		err := rows.Scan(&tb.Condition, &tb.Total)
		if err != nil {
			return nil, err
		}
		// round the total to the second decimal place
		tb.Total = math.Round(tb.Total*100) / 100
		// append the total by condition to the slice
		t = append(t, tb)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// FindTopActive returns the top n active customers in the database by total spent
// total is rounded to the second decimal place.
func (r *CustomersMySQL) FindTopActive(n int) (c []internal.CustomerAmount, err error) {
	// execute the query
	rows, err := r.db.Query(`
        SELECT 
            customers.first_name, 
            customers.last_name, 
            SUM(invoices.total) AS total 
        FROM 
            customers 
        INNER JOIN 
            invoices ON customers.id = invoices.customer_id 
        GROUP BY 
            customers.id 
        ORDER BY 
            total DESC 
        LIMIT ?`,
		n,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var ca internal.CustomerAmount
		// scan the row into the customer amount
		err := rows.Scan(&ca.FirstName, &ca.LastName, &ca.Amount)
		if err != nil {
			return nil, err
		}
		//
		ca.Amount = math.Round(ca.Amount*100) / 100
		// append the customer amount to the slice
		c = append(c, ca)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}
