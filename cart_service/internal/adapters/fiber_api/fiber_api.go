package fiber_api

import (
	"encoding/json"
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

var (
	ErrInvalidId    = errors.New("invalid cart id")
	ErrCartNotFound = errors.New("cart not found")
)

const (
	Add ActionType = iota
	Remove
)

type ActionType uint

type Action struct {
	CartProduct *models.CartProduct `json:"cart_product"`
	Type        ActionType          `json:"type"`
}

type Handler struct {
	app         *fiber.App
	logger      zerolog.Logger
	channels    map[string]chan Action
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func New(logger zerolog.Logger, cartRepo repository.CartRepository, productRepo repository.ProductRepository) *Handler {
	handler := &Handler{
		app:         fiber.New(),
		logger:      logger,
		channels:    make(map[string]chan Action),
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
	handler.RegisterEndpoints()
	handler.app.Use(cors.New())

	return handler
}

func (h *Handler) Listen(addr string) error {
	return h.app.Listen(addr)
}

func (h *Handler) RegisterEndpoints() {
	// TODO: Review these endpoints
	h.app.Get("/cart/:id/ws", h.WebsocketHandler, websocket.New(h.WebsocketManager))
	h.app.Get("/cart/:id", h.GetCart)
	h.app.Post("/cart/:cart_id/products/:product_id", h.AddProduct)
	h.app.Post("/cart/:cart_id/products", h.UpdateProducts)
}

func (h *Handler) WebsocketHandler(ctx *fiber.Ctx) error {
	if !websocket.IsWebSocketUpgrade(ctx) {
		return fiber.ErrUpgradeRequired
	}

	cartId := ctx.Params("id")
	h.logger.Printf("websocket connection for cart %s", cartId)
	h.setupUpdateChannel(cartId)
	return ctx.Next()
}

func (h *Handler) setupUpdateChannel(cartId string) {
	if _, found := h.channels[cartId]; !found {
		h.channels[cartId] = make(chan Action, 10)
	}
}

func (h *Handler) WebsocketManager(c *websocket.Conn) {
	cartId := c.Params("id")

	h.logger.Printf("creating new websocket connection")
	defer h.logger.Printf("closing websocket connection")

	type websocketMessage struct {
		messageType int
		payload     []byte
	}

	var (
		messageType int
		payload     []byte
		err         error
	)

	readerChannel := make(chan websocketMessage)

	reader := func(ch chan<- websocketMessage) {
		for {
			if messageType, payload, err = c.ReadMessage(); err != nil {
				h.logger.Printf("error: %s", err)
				break
			}
			ch <- websocketMessage{messageType, payload}
		}
	}

	go reader(readerChannel)

	defer func() {
		delete(h.channels, cartId)
	}()

	for {
		select {
		case msg := <-readerChannel:
			h.logger.Printf("payload: %s", string(msg.payload))
			if err = c.WriteMessage(messageType, payload); err != nil {
				h.logger.Err(err).Msgf("failed to write message")
				return
			}
		case action := <-h.channels[cartId]:
			h.logger.Printf("update for cart %s", action.CartProduct.CartID)

			payload, err := json.Marshal(action)
			if err != nil {
				h.logger.Err(err).Msgf("failed to notify cart %s", action.CartProduct.CartID)
			}

			if err = c.WriteMessage(websocket.TextMessage, payload); err != nil {
				h.logger.Err(err).Msgf("failed to write message")
				return
			}
		}
	}
}

type UpdateProductsRequestAction string

const (
	add    UpdateProductsRequestAction = "add"
	remove UpdateProductsRequestAction = "remove"
)

// TODO: Validation?
type UpdateProductsRequest struct {
	ProductID string `json:"product_id"`
	Amount    uint   `json:"amount"`
	Action    UpdateProductsRequestAction
}

func (h *Handler) UpdateProducts(ctx *fiber.Ctx) error {
	var body UpdateProductsRequest

	if err := ctx.BodyParser(&body); err != nil {
		return err
	}

	cartId := ctx.Params("cart_id")

	if body.Action == add {
		return h.cartRepo.AddProduct(cartId, body.ProductID, body.Amount)
	} else if body.Action == remove {
		return h.cartRepo.RemoveProduct(cartId, body.ProductID, body.Amount)
	} else {
		return errors.New("invalid action")
	}

	// TODO: Notify
}

func (h *Handler) AddProduct(ctx *fiber.Ctx) error {
	cartId := ctx.Params("cart_id")
	productId := ctx.Params("product_id")
	quantity := uint(1)

	product, err := h.productRepo.GetProduct(productId)
	if err != nil {
		return err
	}

	if err := h.cartRepo.AddProduct(cartId, productId, quantity); err != nil {
		return err
	}

	cp := &models.CartProduct{
		CartID:    cartId,
		ProductID: productId,
		Quantity:  quantity,
		Product:   product,
	}
	h.notify(cp)

	return nil
}

func (h *Handler) notify(cartProduct *models.CartProduct) {
	if cartProduct == nil {
		return
	}

	select {
	case h.channels[cartProduct.CartID] <- Action{Type: Add, CartProduct: cartProduct}:
		// success
		h.logger.Printf("notified cart %s", cartProduct.CartID)
	default:
		// blocked
		h.logger.Printf("ignoring notification for cart %s: channel full", cartProduct.CartID)
	}
}

func (h *Handler) RemoveProduct(ctx *fiber.Ctx) error {
	panic("not implemented")
}

func (h *Handler) GetCart(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ErrInvalidId
	}

	h.logger.Printf("GetCart: %s", id)

	products, err := h.cartRepo.GetCart(id)
	if err != nil {
		return err
	}

	h.logger.Printf("Cart length: %d", len(products))

	if len(products) == 0 {
		return ErrCartNotFound
	}

	cart := &models.Cart{
		ID:       id,
		Products: products,
	}

	return ctx.JSON(cart)
}
