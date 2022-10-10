package fiber_api

import (
	"errors"

	"github.com/fsmiamoto/zcart/cart_service/internal/models"
)

var (
	ErrInvalidId    = errors.New("invalid cart id")
	ErrCartNotFound = errors.New("cart not found")
)

type UpdateProductsRequestAction string

const (
	AddProductAction    UpdateProductsRequestAction = "add"
	RemoveProductAction UpdateProductsRequestAction = "remove"
)

type UpdateProductsRequest struct {
	ProductID string                      `json:"product_id"`
	Quantity  uint                        `json:"quantity"`
	Action    UpdateProductsRequestAction `json:"action"`
}

func (u *UpdateProductsRequest) Validate() error {
	if u.ProductID == "" {
		return errors.New("missing product id")
	}
	if u.Quantity == 0 {
		return errors.New("missing quantity")
	}
	if u.Action == "" {
		return errors.New("missing action")
	}
	return nil
}

// WebSocket
const (
	ProductAddedEvent   = "product_added"
	ProductRemovedEvent = "product_removed"
)

type CartEvent string

type CartEventWebsocketNotification struct {
	CartProduct *models.CartProduct `json:"cart_product"`
	Event       CartEvent           `json:"event"`
}

func updateProductsActionToCartEvent(action UpdateProductsRequestAction) CartEvent {
	if action == RemoveProductAction {
		return ProductRemovedEvent
	}
	return ProductAddedEvent
}
