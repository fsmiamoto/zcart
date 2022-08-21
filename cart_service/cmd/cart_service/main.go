package main

import (
	"database/sql"
	"os"

	"github.com/fsmiamoto/zcart/cart_service/internal/migrations"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository"
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
	logger  = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	devMode = false
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	if os.Getenv("DEV_MODE") == "true" {
		logger.Info().Msgf("Running in Dev Mode")
		devMode = true
		os.Remove(DBFILE)
	}

	db, err := sql.Open("sqlite3", DBFILE)
	fatalIfErr(err)

	if devMode {
		fatalIfErr(migrations.Apply(db))
	}

	cartRepo := repository.NewCartRepository(db)
	productRepo := repository.NewProductRepository(db)

	uihandler := uihandler.New(db, logger, cartRepo, productRepo)

	uihandler.RegisterEndpoints(app)

	fatalIfErr(app.Listen(PORT))
}

func fatalIfErr(err error) {
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}
