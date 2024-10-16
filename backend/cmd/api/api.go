package api

import (
	"database/sql"
	"net/http"

	"github.com/prodanov17/znk/internal/middleware"
	"github.com/prodanov17/znk/internal/services/game"
	"github.com/prodanov17/znk/internal/services/room"
	"github.com/prodanov17/znk/internal/ws"
	"github.com/prodanov17/znk/pkg/logger"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := http.NewServeMux()
	subrouter := http.NewServeMux()
	subrouter.Handle("/api/v1/", http.StripPrefix("/api/v1", router))

	roomRepository := room.NewRepository()
	gameRepository := game.NewRepository()
	gameService := game.NewService(gameRepository)
	roomService := room.NewService(gameService, roomRepository)
	hub := ws.NewHub(roomService)
	wsHandler := ws.NewHandler(hub)
	wsHandler.RegisterRoutes(router)
	go hub.Run()

	// Global middleware stack
	stack := middleware.CreateStack(
		middleware.StripSlashes,
		middleware.Logging,
		middleware.CORS,
	)

	logger.Log.Info("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, stack(subrouter))

}
