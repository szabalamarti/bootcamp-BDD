package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"supermarket/internal"
	"supermarket/platform/web/request"
	"supermarket/platform/web/response"

	"github.com/go-chi/chi/v5"
)

// WarehouseHandler is the handler for the warehouse.
type WarehouseHandler struct {
	// rw is the repository warehouse.
	rw internal.RepositoryWarehouse
}

// NewWarehouseHandler creates a new warehouse handler.
func NewWarehouseHandler(rw internal.RepositoryWarehouse) (wh *WarehouseHandler) {
	wh = &WarehouseHandler{
		rw: rw,
	}
	return
}

// WarehouseJSON is a warehouse in JSON format.
type WarehouseJSON struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

// RequestBodyWarehouseCreate is a request body for creating a warehouse.
type RequestBodyWarehouseCreate struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Telephone string `json:"telephone"`
	Capacity  int    `json:"capacity"`
}

// ReportProductJSON is a report product in JSON format.
type ReportProductJSON struct {
	WarehouseName string `json:"warehouse_name"`
	ProductCount  int    `json:"product_count"`
}

// GetByID returns a warehouse by its id.
func (h *WarehouseHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - path parameter: id
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid id")
			return
		}

		// process
		warehouse, err := h.rw.FindById(id)
		if err != nil {
			switch err {
			case internal.ErrRepositoryWarehouseNotFound:
				response.JSON(w, http.StatusNotFound, "warehouse not found")
				return
			default:
				response.JSON(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}

		// response
		// - serialize warehouse
		data := WarehouseJSON{
			Id:        warehouse.Id,
			Name:      warehouse.Name,
			Address:   warehouse.Address,
			Telephone: warehouse.Telephone,
			Capacity:  warehouse.Capacity,
		}

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"data":    data,
			"message": "warehouse found successfully",
		})
	}
}

// GetAll returns all warehouses.
func (h *WarehouseHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// process
		warehouses, err := h.rw.GetAll()
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize warehouses
		var data []WarehouseJSON
		for _, warehouse := range warehouses {
			data = append(data, WarehouseJSON{
				Id:        warehouse.Id,
				Name:      warehouse.Name,
				Address:   warehouse.Address,
				Telephone: warehouse.Telephone,
				Capacity:  warehouse.Capacity,
			})
		}

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"data":    data,
			"message": "warehouses found successfully",
		})
	}
}

// Create creates a warehouse.
func (h *WarehouseHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - body
		var body RequestBodyWarehouseCreate
		err := request.JSON(r, &body)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid body")
			return
		}

		// process
		warehouse := internal.Warehouse{
			Name:      body.Name,
			Address:   body.Address,
			Telephone: body.Telephone,
			Capacity:  body.Capacity,
		}
		err = h.rw.Save(&warehouse)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize warehouse
		data := WarehouseJSON{
			Id:        warehouse.Id,
			Name:      warehouse.Name,
			Address:   warehouse.Address,
			Telephone: warehouse.Telephone,
			Capacity:  warehouse.Capacity,
		}

		response.JSON(w, http.StatusCreated, map[string]interface{}{
			"data":    data,
			"message": "warehouse created successfully",
		})
	}
}

// GerProductReports returns a reports of products by warehouse.
func (h *WarehouseHandler) GetProductReports() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// request
		// - query param id is a list of warehouse ids
		queryIds, ok := r.URL.Query()["id"]

		var intIds []int
		if ok {
			// Parse the first query parameter as a JSON array
			err := json.Unmarshal([]byte(queryIds[0]), &intIds)
			if err != nil {
				response.JSON(w, http.StatusBadRequest, "invalid id")
				return
			}
		}

		// process
		reports, err := h.rw.ReportProducts(intIds)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// response
		// - serialize reports
		var data []ReportProductJSON
		for _, report := range reports {
			data = append(data, ReportProductJSON{
				WarehouseName: report.WarehouseName,
				ProductCount:  report.ProductCount,
			})
		}

		response.JSON(w, http.StatusOK, map[string]interface{}{
			"data":    data,
			"message": "reports found successfully",
		})
	}
}
