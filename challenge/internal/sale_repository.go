package internal

// RepositorySale is the interface that wraps the basic Sale methods.
type RepositorySale interface {
	// FindAll returns all sales.
	FindAll() (s []Sale, err error)
	// Save saves a sale.
	Save(s *Sale) (err error)
	// FindTopSold returns the top n products sold in the database.
	FindTopSold(n int) (p []ProductSales, err error)
}
