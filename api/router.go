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
)

var Router *fiber.App

func Run(embeddedFiles embed.FS) {

	Router = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
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

	authRoutes := Router.Group("/api", NewAuthenticationMiddleware())
	authRoutes.Get("/clients", GetClients)
	authRoutes.Post("/clients", CreateClient)
	authRoutes.Get("/clients/:id", GetClient)
	authRoutes.Delete("/clients/:id", DeleteClient)
	authRoutes.Get("/config", GetConfig)

	Router.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(assets),
	}))

	go func() {
		err := Router.Listen("0.0.0.0:" + config.Config.HttpsPort)
		if err != nil {
			log.Fatalf("API: Could not start: %s", err)
		}
	}()
}

func Close() {

	err := Router.Shutdown()
	if err != nil {
		log.Fatalf("API: Coult not shut down: %s", err)
	}
}
