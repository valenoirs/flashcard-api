package main

import (
	"log/slog"
	"net/http"

	"github.com/valenoirs/flashcard-api/internal/config"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/command"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/cqrs/query"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/database"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/database/adapter/postgres"
	"github.com/valenoirs/flashcard-api/internal/infrastructure/logger"
	httpRouter "github.com/valenoirs/flashcard-api/internal/transport/http"
	"github.com/valenoirs/flashcard-api/internal/transport/http/middleware"
	"github.com/valenoirs/flashcard-api/internal/usecase/card"
	"github.com/valenoirs/flashcard-api/internal/usecase/deck"
)

type App struct {
	router http.Handler
	logger *slog.Logger
	config *config.Config
}

func NewApp() (*App, func(), error) {
	mux := http.NewServeMux()
	cfg := config.NewConfig()
	log := logger.NewLogger(cfg)

	commandRegistry := command.NewRegistry(log)
	queryRegistry := query.NewRegistry(log)

	postgresWrapper, cleanup, err := database.NewPostgresDatabase(cfg, log)
	if err != nil {
		return nil, nil, err
	}

	// repositories
	cardRepo := postgres.NewCardPostgresAdapter(postgresWrapper)
	deckRepo := postgres.NewDeckPostgresAdapter(postgresWrapper)

	// card handler
	createCardHandler := card.NewCreateCardHandler(cardRepo)
	updateCardHandler := card.NewUpdateCardHandler(cardRepo)
	deleteCardHandler := card.NewDeleteCardHandler(cardRepo)
	getCardListHandler := card.NewGetCardListHandler(cardRepo)
	getCardDetailHandler := card.NewGetCardDetailHandler(cardRepo)

	command.Register(commandRegistry, createCardHandler)
	command.Register(commandRegistry, updateCardHandler)
	command.Register(commandRegistry, deleteCardHandler)
	query.Register(queryRegistry, getCardListHandler)
	query.Register(queryRegistry, getCardDetailHandler)

	// deck handler
	createDeckHandler := deck.NewCreateDeckHandler(deckRepo)
	getDeckListHandler := deck.NewGetDeckListHandler(deckRepo)

	command.Register(commandRegistry, createDeckHandler)
	query.Register(queryRegistry, getDeckListHandler)

	httpRouter.NewHTTPRouter(mux, commandRegistry, queryRegistry)

	// middleware
	mw := middleware.NewManager(cfg, log)

	var globalMiddleware http.Handler = mux

	globalMiddleware = mw.CORS(globalMiddleware)

	return &App{
		router: globalMiddleware,
		logger: log,
		config: cfg,
	}, cleanup, nil
}
