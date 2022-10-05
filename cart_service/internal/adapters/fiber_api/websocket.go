package fiber_api

import (
	// "encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

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
		h.channels[cartId] = make(chan CartEventWebsocketNotification)
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

			if err := c.WriteJSON(action); err != nil {
				h.logger.Err(err).Msgf("failed to write message")
				return
			}

			h.logger.Info().Msgf("notified clients of cart %s", action.CartProduct.CartID)
		}
	}
}
