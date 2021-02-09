package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rafaelcpalmeida/moodz/controllers"
)

type AppConfigurations struct {
	PrefixesFilePath     string
	BusinessSectorAPIURL string
	WebServerHost        string
	WebServerTimeout     int
}

func main() {
	var configs AppConfigurations
	if err := configs.loadConfigs(); err != nil {
		log.Fatalf("Couldn't setup configurations. Error: %v", err)
	}

	srv := controllers.NewWebServer(configs.WebServerHost)
	errChan := make(chan error)
	sigChan := make(chan os.Signal)

	go srv.ListenAndServe(errChan)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Fatalf("cannot start server. Error: %v\n", err)
	default:
	}

	<-sigChan
	log.Println("Attempting to gracefully shut down server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(configs.WebServerTimeout)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func (ac *AppConfigurations) loadConfigs() error {
	ac.WebServerHost = os.Getenv("WEB_SEVER_HOST")
	if ac.WebServerHost == "" {
		return errors.New("env var WEB_SEVER_HOST cannot be empty")
	}

	timeout := os.Getenv("WEB_SEVER_TIMEOUT")
	if timeout == "" {
		ac.WebServerTimeout = 5
	} else {
		timeoutVal, err := strconv.Atoi(timeout)
		if err != nil {
			return errors.New("env var WEB_SEVER_TIMEOUT is invalid")
		}
		ac.WebServerTimeout = timeoutVal
	}

	return nil
}
