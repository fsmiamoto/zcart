package uihandler

import (
	"database/sql"
	_ "embed"
	"errors"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrInvalidId    = errors.New("invalid cart id")
	ErrCartNotFound = errors.New("cart not found")
)

//go:embed migration.sql
var migrations string

type Handler struct {
	db       *sql.DB
	logger   *log.Logger
	channels map[string]chan struct{}
}

func New(db *sql.DB, logger *log.Logger) *Handler {
	return &Handler{
		db:       db,
		logger:   logger,
		channels: make(map[string]chan struct{}),
	}
}

func (h *Handler) RegisterEndpoints(app *fiber.App) {
	app.Get("/cart/:id/ws", h.WebsocketHandler, websocket.New(h.WebsocketManager))
	app.Get("/cart/:id", h.GetCart)
	app.Post("/cart/:id/products/:product_id", h.AddProduct)

	h.applyMigration()
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
		h.channels[cartId] = make(chan struct{}, 10)
	}
}

func (h *Handler) WebsocketManager(c *websocket.Conn) {
	cartId := c.Params("id")

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
		case <-h.channels[cartId]:
			h.logger.Printf("update for cart %s", cartId)
			if err = c.WriteMessage(websocket.TextMessage, []byte("howdy hey")); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func (h *Handler) AddProduct(ctx *fiber.Ctx) error {
	cartId := ctx.Params("id")
	productId := ctx.Params("product_id")
	quantity := 1

	if err := h.addProduct(cartId, productId, quantity); err != nil {
		return err
	}

	h.notify(cartId)

	return nil
}

func (h *Handler) notify(cartId string) {
	select {
	case h.channels[cartId] <- struct{}{}:
		// success
		h.logger.Printf("notified cart %s", cartId)
	default:
		// blocked
		h.logger.Printf("ignoring notification for cart %s: channel full", cartId)
	}
}

func (h *Handler) addProduct(cartId string, productId string, quantity int) error {
	const query = `INSERT INTO cart_products (cart_id,product_id,quantity) VALUES (?,?,?)`

	_, err := h.db.Exec(query, cartId, productId, quantity)

	return err
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
        WHERE cart_id = ?;
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
		return tx.Rollback()
	}

	return tx.Commit()
}
