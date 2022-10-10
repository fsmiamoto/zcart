package fiber_api

import (
	"errors"

	"github.com/fsmiamoto/zcart/cart_service/internal/models"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

// TODO: This is a big ball of mud
// Refactor into the appropriate Application/Domain services

type Handler struct {
	app         *fiber.App
	logger      zerolog.Logger
	channels    map[string]chan CartEventWebsocketNotification
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func New(logger zerolog.Logger, cartRepo repository.CartRepository, productRepo repository.ProductRepository) *Handler {
	handler := &Handler{
		app:         fiber.New(),
		logger:      logger,
		channels:    make(map[string]chan CartEventWebsocketNotification),
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
	handler.app.Use(cors.New())
	handler.RegisterEndpoints()

	return handler
}

func (h *Handler) Listen(addr string) error {
	return h.app.Listen(addr)
}

func (h *Handler) RegisterEndpoints() {
	h.app.Get("/cart/:id/ws", h.WebsocketHandler, websocket.New(h.WebsocketManager))
	h.app.Get("/cart/:id", h.GetCart)
	h.app.Post("/cart/:cart_id/products", h.UpdateProducts)
	h.app.Post("/cart/:cart_id/checkout", h.Checkout)
}

func newError(status int, err error) error {
	return fiber.NewError(status, err.Error())
}

func (h *Handler) Checkout(ctx *fiber.Ctx) error {
	cartId := ctx.Params("cart_id")

	if cartId == "" {
		return newError(fiber.StatusBadRequest, errors.New("missing cart id"))
	}

	h.logger.Info().Msgf("Checkout: %s", cartId)

	if err := h.cartRepo.EmptyCart(cartId); err != nil {
		return err
	}

	return nil
}

func (h *Handler) UpdateProducts(ctx *fiber.Ctx) error {
	// Add missing validation
	var request UpdateProductsRequest

	if err := ctx.BodyParser(&request); err != nil {
		return newError(fiber.StatusBadRequest, err)
	}

	if err := request.Validate(); err != nil {
		return newError(fiber.StatusBadRequest, err)
	}

	cartId := ctx.Params("cart_id")

	product, err := h.productRepo.GetProduct(request.ProductID)
	if err != nil {
		return err
	}

	if err := h.processAction(cartId, request.ProductID, request.Quantity, request.Action); err != nil {
		return err
	}

	cp := &models.CartProduct{
		CartID:    cartId,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
		Product:   product,
	}

	h.notify(cp, request.Action)

	return nil
}

func (h *Handler) GetCart(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ErrInvalidId
	}

	h.logger.Printf("GetCart: %s", id)

	cart, err := h.cartRepo.GetCart(id)
	if err != nil {
		return err
	}

	h.logger.Printf("Cart length: %d", len(cart.Products))

	if cart.Products == nil {
		cart.Products = make([]*models.CartProduct, 0)
	}

	return ctx.JSON(cart)
}

func (h *Handler) processAction(cartId string, productId string, quantity uint, action UpdateProductsRequestAction) error {
	var delta int

	if action == AddProductAction {
		delta = int(quantity)
	} else if action == RemoveProductAction {
		delta = -int(quantity)
	}

	return h.cartRepo.UpdateProductQuantity(cartId, productId, delta)
}

func (h *Handler) notify(cartProduct *models.CartProduct, action UpdateProductsRequestAction) {
	if cartProduct == nil {
		return
	}

	notification := CartEventWebsocketNotification{
		Event:       updateProductsActionToCartEvent(action),
		CartProduct: cartProduct,
	}

	select {
	case h.channels[cartProduct.CartID] <- notification:
		// success
		h.logger.Printf("notified cart %s", cartProduct.CartID)
	default:
		// blocked
		h.logger.Printf("ignoring notification for cart %s: channel full", cartProduct.CartID)
	}
}
