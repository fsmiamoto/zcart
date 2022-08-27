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
}

func (h *Handler) UpdateProducts(ctx *fiber.Ctx) error {
	var request UpdateProductsRequest

	if err := ctx.BodyParser(&request); err != nil {
		return err
	}

	cartId := ctx.Params("cart_id")

	product, err := h.productRepo.GetProduct(request.ProductID)
	if err != nil {
		return err
	}

	if err := h.processAction(cartId, request.ProductID, request.Amount, request.Action); err != nil {
		return err
	}

	cp := &models.CartProduct{
		CartID:    cartId,
		ProductID: request.ProductID,
		Quantity:  request.Amount,
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

	if len(cart.Products) == 0 {
		return ErrCartNotFound
	}

	return ctx.JSON(cart)
}

func (h *Handler) processAction(cartId string, productId string, quantity uint, action UpdateProductsRequestAction) error {
	if action == AddProductAction {
		return h.cartRepo.AddProduct(cartId, productId, quantity)
	} else if action == RemoveProductAction {
		return h.cartRepo.RemoveProduct(cartId, productId, quantity)
	} else {
		return errors.New("invalid action")
	}
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
