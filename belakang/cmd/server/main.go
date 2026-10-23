package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"sekolah/internal/config"
	"sekolah/internal/database"
	"sekolah/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

//go:embed frontend/dist
var frontendFS embed.FS

func main() {
	// Load Config
	if err := config.Load(); err != nil {
		log.Fatal(err)
	}
	// Connect Database
	if err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	app := fiber.New(fiber.Config{
		AppName: config.App.AppName,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(compress.New())
	app.Use(etag.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} ${latency} ${ip} ${bytesSent}B ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Register API Routes DULU (harus sebelum static frontend,
	// supaya /api/... tidak "ketangkep" oleh handler frontend)
	routes.Register(app)

	// Serve file upload (foto siswa, dokumen, dll)
	app.Static("/storage", "./storage", fiber.Static{})

	// Serve frontend Vue (embedded) — HARUS PALING BAWAH
	// karena ini catch-all untuk SPA routing (Vue Router history mode)
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}
	app.Use("/", filesystem.New(filesystem.Config{
		Root:         http.FS(distFS),
		Index:        "index.html",
		NotFoundFile: "index.html", // penting: biar refresh di route Vue (mis. /dashboard/kenaikan) tidak 404
	}))

	log.Printf(
		"[%s] running on :%s",
		config.App.AppName,
		config.App.AppPort,
	)
	log.Fatal(
		app.Listen(":" + config.App.AppPort),
	)
}
