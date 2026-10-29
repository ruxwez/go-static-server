package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/static"
)

// Text types only: images, video and fonts are already compressed. Extensionless
// paths are SPA routes, which serve HTML.
var compressible = map[string]bool{"": true, ".html": true, ".css": true, ".js": true, ".mjs": true, ".json": true, ".svg": true, ".txt": true, ".xml": true}

// newApp serves root. Paths without an extension fall back to index.html so
// client-side routes (React Router, SPA) work; missing assets still return 404.
func newApp(root string) *fiber.App {
	app := fiber.New()
	app.Use(compress.New(compress.Config{
		Next: func(c fiber.Ctx) bool { return !compressible[path.Ext(c.Path())] },
	}))
	app.Use("/", static.New(root, static.Config{
		ByteRange: true, // video and audio seeking needs Range requests
		NotFoundHandler: func(c fiber.Ctx) error {
			if path.Ext(c.Path()) != "" {
				return fiber.ErrNotFound
			}
			// the file server already marked the miss 404; the shell must go out as 200
			return c.Status(fiber.StatusOK).SendFile(filepath.Join(root, "index.html"))
		},
	}))
	return app
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(newApp("static").Listen(":" + port))
}
