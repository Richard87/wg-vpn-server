package api

import (
	"embed"
	"fmt"
	"github.com/Richard87/wg-vpn-server/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var Router *fiber.App

func Run(embeddedFiles embed.FS) {

	Router = fiber.New()
	Router.Use(logger.New())
	Router.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     fmt.Sprintf(fmt.Sprintf("https://localhost:%s, %s", config.Config.HttpsPort, config.Config.HttpsCors)),
		AllowHeaders:     "Origin, Content-Type, Accept, X-Requested-With, Authorization, Access-Control-Allow-Origin",
		AllowMethods:     "GET, HEAD, POST, PUT, OPTIONS, DELETE",
		AllowCredentials: true,
	}))
	assets, err := fs.Sub(embeddedFiles, "ui/build")
	if err != nil {
		log.Fatalf("Could not load UI: %s", err)
	}
	Router.Post("/authenticate", Authenticate)
	Router.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(assets),
	}))

	authRoutes := Router.Group("/api", NewAuthenticationMiddleware())
	authRoutes.Get("/api/clients", GetClients)
	authRoutes.Post("/api/clients", CreateClient)
	authRoutes.Get("/api/clients/:id", GetClient)
	authRoutes.Delete("/api/clients/:id", DeleteClient)
	authRoutes.Get("/api/config", GetConfig)

	go func() {
		err := Router.Listen("0.0.0.0:" + config.Config.HttpsPort)
		if err != nil {
			log.Fatalf("API: Could not start: %s", err)
		}
		log.Println("API: shut down")
	}()

	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGTERM)
	signal.Notify(termSignal, os.Interrupt)
	<-termSignal // Block until we receive our signal.
	log.Println("API: shutting down...")

	err = Router.Shutdown()
	if err != nil {
		log.Fatalf("API: Coult not shut down: %s", err)
	}
}
