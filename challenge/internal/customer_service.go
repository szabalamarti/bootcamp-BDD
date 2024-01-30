package internal

// ServiceCustomer is the interface that wraps the basic methods that a customer service should implement.
type ServiceCustomer interface {
	// FindAll returns all customers
	FindAll() (c []Customer, err error)
	// Save saves a customer
	Save(c *Customer) (err error)
	// FindTotalByCondition returns the aggregated money from invoices by customer condition
	FindTotalByCondition() (t []TotalByCondition, err error)
	// FindTopActive returns the top n active customers in the database by total spent
	FindTopActive(n int) (c []CustomerAmount, err error)
}
