package uihandler

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/fsmiamoto/zcart/cart_service/internal/models"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

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
	logger      zerolog.Logger
	channels    map[string]chan Action
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
}

func New(db *sql.DB, logger zerolog.Logger, cartRepo repository.CartRepository, productRepo repository.ProductRepository) *Handler {
	return &Handler{
		logger:      logger,
		channels:    make(map[string]chan Action),
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (h *Handler) RegisterEndpoints(app *fiber.App) {
	app.Get("/cart/:id/ws", h.WebsocketHandler, websocket.New(h.WebsocketManager))
	app.Get("/cart/:id", h.GetCart)
	app.Post("/cart/:id/products/:product_id", h.AddProduct)
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

func (h *Handler) AddProduct(ctx *fiber.Ctx) error {
	cartId := ctx.Params("id")
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
