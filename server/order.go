package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aeilang/nice/auth"
	"github.com/aeilang/nice/db/store"
	"github.com/aeilang/nice/utils"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func (s *Server) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusNetworkAuthenticationRequired, nil)
	}

	var cart CardCheckoutPayload
	if err := utils.ParseJSON(r, &cart); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(cart); err != nil {
		errs := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errs))
		return
	}

	productIds, err := getCartItemsIDs(cart.Items)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	products, err := s.Querier.GetProductsByIDs(r.Context(), productIds)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	orderID, totalPrice, err := s.createOrder(products, cart.Items, userID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"total_price": totalPrice.Int,
		"order_id":    orderID,
	})
}

type CardCheckoutPayload struct {
	Items []CardCheckoutItem `json:"items" validate:"required"`
}

type CardCheckoutItem struct {
	ProductID int32 `json:"product_id" validate:"required"`
	Quantity  int32 `json:"quantity" validate:"required,gt=0"`
}

func getCartItemsIDs(items []CardCheckoutItem) ([]int32, error) {
	productIds := make([]int32, len(items))
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product %d", item.ProductID)
		}

		productIds[i] = item.ProductID
	}

	return productIds, nil
}

func (s *Server) createOrder(products []store.Product, carItems []CardCheckoutItem, useID int) (int, pgtype.Numeric, error) {
	productMap := make(map[int32]store.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	if err := checkIfCardIsInStock(carItems, productMap); err != nil {
		return 0, pgtype.Numeric{}, err
	}

	totalPrice := calculateTotalPrice(carItems, productMap)

	tx, err := s.Pool.Begin(context.Background())
	if err != nil {
		return 0, pgtype.Numeric{}, err
	}

	defer tx.Rollback(context.Background())
	query := store.New(tx)

	// reduce the quantity of products in the store
	for _, item := range carItems {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		err := query.UpdateProduct(context.Background(), store.UpdateProductParams{
			Name:        product.Name,
			Price:       product.Price,
			Image:       product.Image,
			Description: product.Description,
			Quantity:    product.Quantity,
			ID:          product.ID,
		})
		if err != nil {
			return 0, pgtype.Numeric{}, err
		}
	}

	// create order record
	order, err := query.CreateOrder(context.Background(), store.CreateOrderParams{
		UserID:  int32(useID),
		Total:   totalPrice,
		Status:  "pending",
		Address: "some address",
	})

	if err != nil {
		return 0, pgtype.Numeric{}, err
	}

	// create order the items records
	for _, item := range carItems {

		query.CreateOrderItem(context.Background(), store.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	err = tx.Commit(context.Background())

	return int(order.ID), totalPrice, err
}

func checkIfCardIsInStock(cartItems []CardCheckoutItem, products map[int32]store.Product) error {
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	for _, item := range cartItems {
		product, ok := products[item.ProductID]
		if !ok {
			return fmt.Errorf("product %d is not available in the stock", item.ProductID)
		}

		if product.Quantity < item.Quantity {
			return fmt.Errorf("product %s is not avaliable in the qunatity", product.Name)
		}
	}

	return nil
}

func calculateTotalPrice(cartItems []CardCheckoutItem, products map[int32]store.Product) pgtype.Numeric {
	total := decimal.NewFromInt(0)

	for _, item := range cartItems {
		product := products[item.ProductID]
		total = total.Add(utils.ToDecimal(product.Price))
	}

	return utils.ToNumeric(total)
}
