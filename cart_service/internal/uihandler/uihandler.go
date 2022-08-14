package uihandler

import (
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrInvalidId       = errors.New("invalid cart id")
	ErrCartNotFound    = errors.New("cart not found")
	ErrProductNotFound = errors.New("product not found")
)

//go:embed migration.sql
var migrations string

const (
	Add ActionType = iota
	Remove
)

type ActionType uint

type Action struct {
	CartProduct *CartProduct `json:"cart_product"`
	Type        ActionType   `json:"type"`
}

type Handler struct {
	db       *sql.DB
	logger   *log.Logger
	channels map[string]chan Action
}

func New(db *sql.DB, logger *log.Logger) *Handler {
	return &Handler{
		db:       db,
		logger:   logger,
		channels: make(map[string]chan Action),
	}
}

func (h *Handler) RegisterEndpoints(app *fiber.App) error {
	app.Get("/cart/:id/ws", h.WebsocketHandler, websocket.New(h.WebsocketManager))
	app.Get("/cart/:id", h.GetCart)
	app.Post("/cart/:id/products/:product_id", h.AddProduct)

	if err := h.applyMigration(); err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}
	return nil
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
				log.Println("write:", err)
				return
			}
		case action := <-h.channels[cartId]:
			h.logger.Printf("update for cart %s", action.CartProduct.CartID)

			payload, err := json.Marshal(action)
			if err != nil {
				log.Printf("failed to notify cart %s: %s", action.CartProduct.CartID, err)
			}

			if err = c.WriteMessage(websocket.TextMessage, payload); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func (h *Handler) AddProduct(ctx *fiber.Ctx) error {
	cartId := ctx.Params("id")
	productId := ctx.Params("product_id")
	quantity := uint(1)

	product, err := h.getProduct(productId)
	if err != nil {
		return err
	}

	if err := h.addProduct(cartId, productId, quantity); err != nil {
		return err
	}

	cp := &CartProduct{
		CartID:    cartId,
		ProductID: productId,
		Quantity:  quantity,
		Product:   product,
	}
	h.notify(cp)

	return nil
}

func (h *Handler) notify(cartProduct *CartProduct) {
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

func (h *Handler) addProduct(cartId string, productId string, quantity uint) error {
	const query = `INSERT INTO cart_products (cart_id,product_id,quantity) VALUES (?,?,?)`

	_, err := h.db.Exec(query, cartId, productId, quantity)

	return err
}

func (h *Handler) getProduct(productId string) (Product, error) {
	const query = `SELECT id, name, price, description, image_url FROM products WHERE id = ?`
	var product Product

	row := h.db.QueryRow(query, productId)

	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Description, &product.ImageURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return product, ErrProductNotFound
		}
		return product, err
	}

	return product, nil
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

	products, err := h.getCart(id)
	if err != nil {
		return err
	}

	h.logger.Printf("Cart length: %d", len(products))

	if len(products) == 0 {
		return ErrCartNotFound
	}

	cart := &Cart{
		ID:       id,
		Products: products,
	}

	return ctx.JSON(cart)
}

type Cart struct {
	ID       string         `json:"id"`
	Products []*CartProduct `json:"products"`
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    *string `json:"image_url"`
}

type CartProduct struct {
	CartID    string  `json:"cart_id"`
	ProductID string  `json:"product_id"`
	Quantity  uint    `json:"quantity"`
	Product   Product `json:"product"`
}

func (h *Handler) getCart(cartId string) ([]*CartProduct, error) {
	const query = `
    SELECT 
        cp.cart_id, cp.product_id, cp.quantity, 
        p.name, p.price, p.id, p.description, p.image_url
        FROM cart_products cp 
        JOIN products p ON cp.product_id = p.id 
        WHERE cart_id = ?
        ORDER BY cp.updated_at;
    `

	var result []*CartProduct

	rows, err := h.db.Query(query, cartId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		cp := &CartProduct{}
		if err := rows.Scan(
			&cp.CartID, &cp.ProductID, &cp.Quantity, &cp.Product.Name,
			&cp.Product.Price, &cp.Product.ID, &cp.Product.Description, &cp.Product.ImageURL,
		); err != nil {
			return nil, err
		}
		result = append(result, cp)
	}

	return result, nil
}

func (h *Handler) applyMigration() error {
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(migrations)
	if err != nil {
		defer tx.Rollback()
		return err
	}

	return tx.Commit()
}
