package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	address string
	port    int
	server  *http.Server
	router  *gin.Engine
}

func NewWebServer(
	address string,
	port int,
) Server {
	httpserver := Server{
		address: address,
		port:    port,
	}

	router := httpserver.setupRouter()

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: router,
	}

	httpserver.server = srv
	httpserver.router = router

	return httpserver
}

func (s Server) ListenAndServe(c chan error) {
	log.Printf("Ready to handle requests on %s:%d\n", s.address, s.port)
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
		c.JSON(http.StatusOK, gin.H{
			"message": "Hi, visitor!",
		})
	})

	return r
}
