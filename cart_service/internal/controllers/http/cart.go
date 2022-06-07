package http

import (
	"strconv"

	"github.com/fsmiamoto/zcart/cart_service/internal/usecases"
	"github.com/gofiber/fiber/v2"
)

type cartHandler struct {
	cartUsecase usecases.Cart
}

func RegisterCartHandlers(app *fiber.App, u usecases.Cart) {
	handler := &cartHandler{
		cartUsecase: u,
	}
	app.Get("cart/:id/products", handler.ListCartProducts)
}

func (h *cartHandler) ListCartProducts(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return err
	}
	cart, err := h.cartUsecase.GetByID(id)
	if err != nil {
		return err
	}
	cartProds, err := h.cartUsecase.ListProducts(cart)
	if err != nil {
		return err
	}
	return ctx.JSON(cartProds)
}
