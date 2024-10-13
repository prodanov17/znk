package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/prodanov17/znk/internal/middleware"
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

	hub := ws.NewHub()
	wsHandler := ws.NewHandler(hub)
	wsHandler.RegisterRoutes(router)
	go hub.Run()

	// Global middleware stack
	stack := middleware.CreateStack(
		middleware.StripSlashes,
		middleware.Logging,
		middleware.CORS,
	)

	log.Println("Starting server on", s.addr)

	return http.ListenAndServe(s.addr, stack(subrouter))

}
