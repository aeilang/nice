package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aeilang/nice/db/store"
	"github.com/aeilang/nice/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func (s *Server) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	str := chi.URLParam(r, "id")
	if str == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing path variable"))
		return
	}

	id, err := strconv.Atoi(str)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid path variable"))
		return
	}

	product, err := s.Querier.GetProductByID(r.Context(), int32(id))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, product)
}

func (s *Server) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := s.Querier.GetProducts(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, products)
}

func (s *Server) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	var product CreateProductPayload
	if err := utils.ParseJSON(r, &product); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(product); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errs))
		return
	}

	price, err := decimal.NewFromString(product.Price)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	p, err := s.Querier.CreateProduct(r.Context(), store.CreateProductParams{
		Name:        product.Name,
		Description: product.Description,
		Image:       product.Image,
		Price:       utils.ToNumeric(price),
		Quantity:    product.Quantity,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, p)
}

type CreateProductPayload struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Price       string `json:"price" validate:"required"`
	Quantity    int32  `json:"quantity" validate:"required"`
}
