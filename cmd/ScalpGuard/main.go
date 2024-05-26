package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/joelpm3/ScalpGuard/pkg/server"// Adjust this import path to where your server package is located.
	"github.com/joelpm3/ScalpGuard/pkg/config"// Adjust this import path to where your config package is located.
)

func main() {
	// Initialize server configuration.
	// This could also be done through environment variables or a config file.

	serverConfig, err := config.LoadServerConfig("config/server.yaml")
	if err != nil {
		log.Fatalf("Failed to load server config: %v", err)
	}

	browserConfig, err := config.LoadBrowserConfig("config/browser.yaml")
	if err != nil {
		log.Fatalf("Failed to load browser config: %v", err)
	}
		

	config := &server.Config{
		Host:    "0.0.0.0",
		Port:    80,     // HTTP port, not used in the simplified version but could be if you decide to support HTTP.
		TlsPort: 443,    // HTTPS port
	}

	// Create a new server instance.
	srv := server.New(config, &serverConfig, &browserConfig)

	// Create a context that is canceled on system interrupt or SIGTERM signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Run the server.
	if err := srv.Run(ctx); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

	// Wait for context cancellation (which happens on interrupt or SIGTERM).
	<-ctx.Done()

	// Attempt to gracefully shutdown the server here, if needed.
	log.Println("Shutting down server...")
}
