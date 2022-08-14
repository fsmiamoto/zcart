package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/fsmiamoto/zcart/cart_service/internal/uihandler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	os.Remove("./zcart.db")

	db, err := sql.Open("sqlite3", "./zcart.db")
	fatalIfErr(err)

	uihandler := uihandler.New(db, log.New(os.Stdout, "", 0))

	fatalIfErr(uihandler.RegisterEndpoints(app))

	log.Fatal(app.Listen(":3333"))
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
