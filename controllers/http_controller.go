package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Address                  string
	server                   *http.Server
	router                   *gin.Engine
}

func NewWebServer(
	address string,
) Server {
	httpserver := Server{
		Address:                  address,
	}

	router := httpserver.setupRouter()

	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}

	httpserver.server = srv
	httpserver.router = router

	return httpserver
}

func (s Server) ListenAndServe(c chan error) {
	log.Printf("Ready to handle requests on %s\n", s.Address)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("Got err. %s\n", err.Error())
		c <- err
	}
}

func (s Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s Server) setupRouter() *gin.Engine {
	log.Println("Registering HTTP endpoints")
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Hi, visitor.")
	})

	return r
}
