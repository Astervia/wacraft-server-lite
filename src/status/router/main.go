package status_router

import (
	auth_middleware "github.com/Astervia/wacraft-server/src/auth/middleware"
	status_handler "github.com/Astervia/wacraft-server/src/status/handler"
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	group := app.Group("/status")

	mainRoutes(group)
	whatsappRoutes(group)
}

func mainRoutes(group fiber.Router) {
	group.Get("", auth_middleware.UserMiddleware, status_handler.Get)
	group.Get("/count", auth_middleware.UserMiddleware, status_handler.Count)
	group.Get("/content/like/:likeText", auth_middleware.UserMiddleware, status_handler.ContentLike)
	group.Get("/content/:keyName/like/:likeText", auth_middleware.UserMiddleware, status_handler.ContentKeyLike)
}
