package loader

import (
	"app/internal"
	"encoding/json"
	"os"
)

type InvoiceLoaderJSON struct {
	rp       internal.RepositoryInvoice
	filepath string
}

func NewInvoiceLoaderJSON(rp internal.RepositoryInvoice, filepath string) *InvoiceLoaderJSON {
	return &InvoiceLoaderJSON{rp, filepath}
}

type InvoiceJSON struct {
	Id         int     `json:"id"`
	Datetime   string  `json:"datetime"`
	Total      float64 `json:"total"`
	CustomerId int     `json:"customer_id"`
}

// Load invoices from JSON file.
func (l *InvoiceLoaderJSON) Load() (c []internal.Invoice, err error) {
	var cs []InvoiceJSON
	// open the file
	f, err := os.Open(l.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// the file has a JSON array, so we need to decode it into a slice
	err = json.NewDecoder(f).Decode(&cs)
	if err != nil {
		return nil, err
	}

	// iterate over the slice and append the invoices to the slice
	for _, v := range cs {
		c = append(c, internal.Invoice{
			Id: v.Id,
			InvoiceAttributes: internal.InvoiceAttributes{
				Datetime:   v.Datetime,
				Total:      v.Total,
				CustomerId: v.CustomerId,
			},
		})
	}

	return
}

// Migrate invoices to the repository.
func (l *InvoiceLoaderJSON) Migrate() (err error) {
	// load the invoices
	c, err := l.Load()
	if err != nil {
		return err
	}

	// iterate over the slice and append the invoices to the repository
	for _, v := range c {
		err = l.rp.Save(&v)
		if err != nil {
			return err
		}
	}

	return
}
