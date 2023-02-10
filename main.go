package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/magesh1/go-mongo/db"
	"github.com/magesh1/go-mongo/middleware"
	"github.com/magesh1/go-mongo/routes"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("app")       // name of config file (without extension)
	viper.AddConfigPath("./config/") // path to look for the config file in
	viper.SetConfigType("yaml")      // REQUIRED if the config file does not have the extension in the name
	err := viper.ReadInConfig()      // Find and read the config file
	if err != nil {                  // Handle errors reading the config file
		log.Error().Msg("error reading config file: " + err.Error())
	}

}

func main() {

	// connecting to the database
	db.DbInit()

	router := gin.Default()

	// setup rate limit
	router.Use(middleware.RateLimit())

	routes.Routes(router)

	srv := &http.Server{
		Addr:              viper.GetString("server.port"),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %s\n", err)
		}
	}()

	fmt.Println("Server is running on port 8080")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Error shutting down server: %s\n", err)
	}

	fmt.Println("Server shutdown successfully")
}
