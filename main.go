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
	WebServerPort        int
	WebServerTimeout     int
}

func main() {
	var configs AppConfigurations
	if err := configs.loadConfigs(); err != nil {
		log.Fatalf("Couldn't setup configurations. Error: %v", err)
	}

	srv := controllers.NewWebServer(configs.WebServerHost, configs.WebServerPort)
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
	ac.WebServerHost = os.Getenv("HOST")

	webServerPortValue, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		return errors.New("env var PORT is not a number")
	}

	ac.WebServerPort = webServerPortValue
	if ac.WebServerPort == 0 {
		return errors.New("env var PORT cannot be zero")
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
