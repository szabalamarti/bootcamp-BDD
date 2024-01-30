package internal

// ServiceCustomer is the interface that wraps the basic methods that a customer service should implement.
type ServiceCustomer interface {
	// FindAll returns all customers
	FindAll() (c []Customer, err error)
	// Save saves a customer
	Save(c *Customer) (err error)
	// GetTotalByCondition returns the aggregated money from invoices by customer condition
	GetTotalByCondition() (t []TotalByCondition, err error)
}
