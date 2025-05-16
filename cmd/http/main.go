package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"clean-arch/cmd/http/app"

	_ "github.com/google/subcommands"
)

//	@title			Babyplug Clean Arch API
//	@version		0.1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Babyplug
//	@contact.email	poramin.lertudom@gmail.com

//	@license.name	No License

//	@host		localhost:8080
//	@BasePath	/v1
//	@schemes	http

// @securityDefinitions.basic	BearerAuth
//
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	ctx := context.Background()

	app, err := app.InitializeApplication(ctx)
	defer app.MongoClient.Disconnect(ctx)
	if err != nil {
		log.Fatalf("Failed to initialized application: %v", err)
	}

	// Background user logger
	stopCh := make(chan struct{})
	app.StartBackgroundProcess(stopCh)

	// Start server in goroutine
	go func() {
		log.Println("Server started on " + app.Config.Port)
		if err := app.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server listen error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	close(stopCh)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown: ", err)
	}

	log.Println("Servers gracefully stopped")

	// catching ctx.Done(). timeout of 10 seconds.
	<-ctx.Done()
	log.Println("timeout of 10 seconds.")
	log.Println("Server exiting")
}
