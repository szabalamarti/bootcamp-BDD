package service

import "app/internal"

// NewCustomersDefault creates new default service for customer entity.
func NewCustomersDefault(rp internal.RepositoryCustomer) *CustomersDefault {
	return &CustomersDefault{rp}
}

// CustomersDefault is the default service implementation for customer entity.
type CustomersDefault struct {
	// rp is the repository for customer entity.
	rp internal.RepositoryCustomer
}

// FindAll returns all customers.
func (s *CustomersDefault) FindAll() (c []internal.Customer, err error) {
	c, err = s.rp.FindAll()
	return
}

// Save saves the customer.
func (s *CustomersDefault) Save(c *internal.Customer) (err error) {
	err = s.rp.Save(c)
	return
}

// FindTotalByCondition returns the aggregated money from invoices by customer condition.
func (s *CustomersDefault) FindTotalByCondition() (t []internal.TotalByCondition, err error) {
	t, err = s.rp.FindTotalByCondition()
	return
}

// FindTopActive returns the top n active customers in the database by total spent.
func (s *CustomersDefault) FindTopActive(n int) (c []internal.CustomerAmount, err error) {
	c, err = s.rp.FindTopActive(n)
	return
}
