package middleware

import (
	"fmt"
	"strings"

	"klinik-pkp-api/internal/api/v1/viewer"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var viewerChan = make(chan viewer.Viewer, 100)

func StartViewerWorker(db *gorm.DB) {
	go func() {
		for v := range viewerChan {
			db.Create(&v)
			fmt.Println("success create viewer")
		}
	}()
}

func ViewerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// skip admin
		role, ok := c.Locals("role").(string)
		if ok && role == "admin" {
			fmt.Printf("test viewer middleware %s\n", role)
			return c.Next()
		}

		// skip internal/admin APIs if needed
		if strings.Contains(c.Path(), "/api/v1") {
			fmt.Println("got here")
			viewerChan <- viewer.Viewer{
				Path:   c.Path(),
				Method: c.Method(),
			}
		}

		return c.Next()
	}
}
