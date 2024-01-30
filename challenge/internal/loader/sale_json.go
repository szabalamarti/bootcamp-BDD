package loader

import (
	"app/internal"
	"encoding/json"
	"os"
)

type SaleLoaderJSON struct {
	rp       internal.RepositorySale
	filepath string
}

func NewSaleLoaderJSON(rp internal.RepositorySale, filepath string) *SaleLoaderJSON {
	return &SaleLoaderJSON{rp, filepath}
}

type SaleJSON struct {
	Id        int `json:"id"`
	Quantity  int `json:"quantity"`
	ProductId int `json:"product_id"`
	InvoiceId int `json:"invoice_id"`
}

// Load sales from JSON file.
func (l *SaleLoaderJSON) Load() (c []internal.Sale, err error) {
	var cs []SaleJSON
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

	// iterate over the slice and append the sales to the slice
	for _, v := range cs {
		c = append(c, internal.Sale{
			Id: v.Id,
			SaleAttributes: internal.SaleAttributes{
				Quantity:  v.Quantity,
				ProductId: v.ProductId,
				InvoiceId: v.InvoiceId,
			},
		})
	}

	return
}

// Migrate sales to the repository.
func (l *SaleLoaderJSON) Migrate() (err error) {
	// load the sales
	c, err := l.Load()
	if err != nil {
		return err
	}

	// migrate the sales
	for _, v := range c {
		err = l.rp.Save(&v)
		if err != nil {
			return err
		}
	}

	return
}
