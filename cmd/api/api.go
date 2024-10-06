package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/prodanov17/znk/internal/middleware"
	"github.com/prodanov17/znk/internal/services/user"
	"github.com/prodanov17/znk/internal/ws"
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

	userStore := user.NewStore(s.db)
	userService := user.NewService(userStore)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRoutes(router)

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	wsHandler.RegisterRoutes(router)
	go hub.Run()

	// Global middleware stack
	stack := middleware.CreateStack(
		middleware.StripSlashes,
		middleware.Logging,
	)

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, stack(subrouter))

}
