package internal

type InvoiceLoader interface {
	Load() (c []Invoice, err error)
	Migrate() (err error)
}
