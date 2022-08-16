package main

import (
	"database/sql"
	"os"

	"github.com/fsmiamoto/zcart/cart_service/internal/uihandler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// TODO: Make this configurable
const (
	PORT   = ":3333"
	DBFILE = "./zcart.db"
)

var (
	logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	if os.Getenv("DEV_MODE") == "true" {
		logger.Info().Msgf("Running in DEV_MODE")
		os.Remove(DBFILE)
	}

	db, err := sql.Open("sqlite3", DBFILE)
	fatalIfErr(err)

	uihandler := uihandler.New(db, logger)

	fatalIfErr(uihandler.RegisterEndpoints(app))

	fatalIfErr(app.Listen(PORT))
}

func fatalIfErr(err error) {
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}
