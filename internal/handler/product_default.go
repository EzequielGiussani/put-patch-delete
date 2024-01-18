package handler

import (
	"app/internal"
	"app/platform/web/request"
	"app/platform/web/response"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DefaultProduct struct {
	sv internal.ProductService
}

type BodyRequestProductJSON struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type BodyResponseProductJSON struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func NewDefaultProducts(sv internal.ProductService) *DefaultProduct {
	return &DefaultProduct{
		sv: sv,
	}
}

func ValidateKeyExistance(mp map[string]any, keys ...string) error {
	for _, key := range keys {
		if _, ok := mp[key]; !ok {
			return fmt.Errorf("key %s does not exist", key)
		}
	}
	return nil
}

func (d *DefaultProduct) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var mp map[string]any

		requestBody, err := io.ReadAll(r.Body)

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body"))
			return
		}

		if err := json.Unmarshal(requestBody, &mp); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body" + err.Error()))
			return
		}

		if err := ValidateKeyExistance(mp, "name", "quantity", "code_value", "expiration", "price"); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body " + err.Error()))
			return
		}

		var body BodyRequestProductJSON

		if err := json.Unmarshal(requestBody, &body); err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid body" + err.Error()))
			return
		}

		product := internal.Product{
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		if err := d.sv.Save(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrFieldRequired), errors.Is(err, internal.ErrProductCodeAlreadyExists), errors.Is(err, internal.ErrFieldFormat):
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("invalid body: " + err.Error()))
			default:
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("internal server error"))
			}
			return
		}

		data := BodyResponseProductJSON{
			ID:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"Message": "Product created successfully",
			"data":    data,
		})

	}

}

func (d *DefaultProduct) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid id"))
			return
		}

		product, err := d.sv.GetById(id)

		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductID):
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("product with the provided id not found"))
				return
			}
		}

		response := BodyResponseProductJSON{
			ID:          product.ID,
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"Message": "Product found successfully",
			"data":    response,
		})

	}
}

func (d *DefaultProduct) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		bytes, err := io.ReadAll(r.Body)

		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		var bodyMap map[string]any
		if err := json.Unmarshal(bytes, &bodyMap); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		if err := ValidateKeyExistance(bodyMap, "name", "quantity", "code_value", "expiration", "price"); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		var body BodyRequestProductJSON
		if err := json.Unmarshal(bytes, &body); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		product := internal.Product{
			ID:          id,
			Name:        body.Name,
			Quantity:    body.Quantity,
			CodeValue:   body.CodeValue,
			IsPublished: body.IsPublished,
			Expiration:  body.Expiration,
			Price:       body.Price,
		}

		if err := d.sv.Update(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product with the provided id not found")
				return
			case errors.Is(err, internal.ErrFieldRequired), errors.Is(err, internal.ErrProductCodeAlreadyExists), errors.Is(err, internal.ErrFieldFormat):
				response.Text(w, http.StatusBadRequest, "invalid body: "+err.Error())
				return
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}
	}
}

func (d *DefaultProduct) UpdatePartial() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
			return
		}

		product, err := d.sv.GetById(id)
		if err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product with the provided id not found")
				return
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
				return
			}
		}

		reqBody := BodyRequestProductJSON{
			Name:        product.Name,
			Quantity:    product.Quantity,
			CodeValue:   product.CodeValue,
			IsPublished: product.IsPublished,
			Expiration:  product.Expiration,
			Price:       product.Price,
		}

		if err := request.JSON(r, &reqBody); err != nil {
			response.Text(w, http.StatusBadRequest, "invalid body")
			return
		}

		product = internal.Product{
			ID:          id,
			Name:        reqBody.Name,
			Quantity:    reqBody.Quantity,
			CodeValue:   reqBody.CodeValue,
			IsPublished: reqBody.IsPublished,
			Expiration:  reqBody.Expiration,
			Price:       reqBody.Price,
		}

		if err := d.sv.Update(&product); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product with the provided id not found")
			case errors.Is(err, internal.ErrFieldRequired), errors.Is(err, internal.ErrProductCodeAlreadyExists), errors.Is(err, internal.ErrFieldFormat):
				response.Text(w, http.StatusBadRequest, "invalid body: "+err.Error())
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		data := BodyResponseProductJSON{
			ID:          id,
			Name:        reqBody.Name,
			Quantity:    reqBody.Quantity,
			CodeValue:   reqBody.CodeValue,
			IsPublished: reqBody.IsPublished,
			Expiration:  reqBody.Expiration,
			Price:       reqBody.Price,
		}

		response.JSON(w, http.StatusOK, map[string]any{
			"Message": "Product updated successfully",
			"data":    data,
		})

	}
}

func (d *DefaultProduct) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(chi.URLParam(r, "id"))

		if err != nil {
			response.Text(w, http.StatusBadRequest, "invalid id")
		}

		if err := d.sv.Delete(id); err != nil {
			switch {
			case errors.Is(err, internal.ErrProductNotFound):
				response.Text(w, http.StatusNotFound, "product with the provided id not found")
			default:
				response.Text(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		response.Text(w, http.StatusOK, "Product deleted successfully")
	}
}
