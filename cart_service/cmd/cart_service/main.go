package main

import (
	"database/sql"
	"os"

	fiberApi "github.com/fsmiamoto/zcart/cart_service/internal/adapters/fiber_api"
	"github.com/fsmiamoto/zcart/cart_service/internal/migrations"
	"github.com/fsmiamoto/zcart/cart_service/internal/repository/sqlite"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
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

	cartRepo := sqlite.NewCartRepository(db)
	productRepo := sqlite.NewProductRepository(db)

	api := fiberApi.New(logger, cartRepo, productRepo)

	fatalIfErr(api.Listen(PORT))
}

func fatalIfErr(err error) {
	if err != nil {
		logger.Fatal().Err(err).Msg("")
	}
}
