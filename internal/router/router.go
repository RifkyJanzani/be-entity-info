package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/rifky/be-entity-info/internal/handler"
)

// validKinds is the set of allowed entity kinds.
var validKinds = map[string]bool{
	"Vehicle":    true,
	"IoT Device": true,
	"Facility":   true,
	"Asset":      true,
}

// validStatuses is the set of allowed entity statuses.
var validStatuses = map[string]bool{
	"Active":      true,
	"Idle":        true,
	"Maintenance": true,
	"Offline":     true,
}

// NewRouter creates and configures the Gin engine with middleware and routes.
func NewRouter(h *handler.EntityHandler, allowedOrigins []string, appEnv string) *gin.Engine {
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Register custom validators for entity kind and status
	registerCustomValidators()

	// CORS middleware
	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	if appEnv == "production" {
		corsConfig.AllowOrigins = allowedOrigins
	} else {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
	}

	r.Use(cors.New(corsConfig))

	// API routes
	v1 := r.Group("/api/v1")
	{
		entities := v1.Group("/entities")
		{
			entities.GET("", h.List)
			entities.GET("/metrics", h.GetMetrics) // Must be before /:id
			entities.GET("/:id", h.GetByID)
			entities.POST("", h.Create)
			entities.PUT("/:id", h.Update)
			entities.DELETE("/:id", h.Delete)
		}
	}

	return r
}

// registerCustomValidators registers custom validation tags with Gin's validator.
func registerCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("entity_kind", func(fl validator.FieldLevel) bool {
			return validKinds[fl.Field().String()]
		})
		v.RegisterValidation("entity_status", func(fl validator.FieldLevel) bool {
			return validStatuses[fl.Field().String()]
		})
	}
}
