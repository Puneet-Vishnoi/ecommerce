package router

import (
	"ecommerce-project/auth"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(*gin.Context)
}
type routes struct {
	router *gin.Engine
}

type Routes []Route

func (r routes) EcommerceUser(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSEMiddleware())
	for _, route := range userRoutes {
		switch route.Method {
			case "GET":
				orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
			case "POST":
				orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
			case "OPTIONS":
				orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
			case "PUT":
				orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
			case "DELETE":
				orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
			default:
				orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
					c.JSON(200, gin.H{
						"result": "Specify a valid http method with this route.",
					})
				})
		}
	}
}


func (r routes) EcommerceProduct(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSEMiddleware())
	for _, route := range productRoutes {
		switch route.Method {
			case "GET":
				orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
			case "POST":
				orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
			case "OPTIONS":
				orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
			case "PUT":
				orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
			case "DELETE":
				orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
			default:
				orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
					c.JSON(200, gin.H{
						"result": "Specify a valid http method with this route.",
					})
				})
		}
	}
}

func (r routes) EcommerceGlobalProductRoutes(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce-product")
	orderRouteGrouping.Use(CORSEMiddleware())
	for _, route := range productGlobalRoutes {
		switch route.Method {
			case "GET":
				orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
			case "POST":
				orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
			case "OPTIONS":
				orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
			case "PUT":
				orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
			case "DELETE":
				orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
			default:
				orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
					c.JSON(200, gin.H{
						"result": "Specify a valid http method with this route.",
					})
				})
		}
	}
}

func (r routes) EcommerceAuthUser(rg *gin.RouterGroup) {
	orderRouteGrouping := rg.Group("/ecommerce")
	orderRouteGrouping.Use(CORSEMiddleware())
	for _, route := range userAuthRoutes {
		switch route.Method {
			case "GET":
				orderRouteGrouping.GET(route.Pattern, route.HandlerFunc)
			case "POST":
				orderRouteGrouping.POST(route.Pattern, route.HandlerFunc)
			case "OPTIONS":
				orderRouteGrouping.OPTIONS(route.Pattern, route.HandlerFunc)
			case "PUT":
				orderRouteGrouping.PUT(route.Pattern, route.HandlerFunc)
			case "DELETE":
				orderRouteGrouping.DELETE(route.Pattern, route.HandlerFunc)
			default:
				orderRouteGrouping.GET(route.Pattern, func(c *gin.Context) {
					c.JSON(200, gin.H{
						"result": "Specify a valid http method with this route.",
					})
				})
		}
	}
}

// append routes with versions
func ClientRoutes() {
	r := routes{
		router: gin.Default(),
	}

	v1 := r.router.Group(os.Getenv("API_VERSION"))
	r.EcommerceUser(v1)
	r.EcommerceGlobalProductRoutes(v1)

	v1.Use(auth.Auth())
	r.EcommerceProduct(v1)
	r.EcommerceAuthUser(v1)

	if err := r.router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Printf("Failed to run server: %v", err)
	}
}

// Middlewares
func CORSEMiddleware() gin.HandlerFunc {
	return func (c *gin.Context){
		c.Writer.Header().Set("Content-type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS"{
			c.Status(http.StatusOK)
		}
	}
}