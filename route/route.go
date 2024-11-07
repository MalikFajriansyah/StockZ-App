package route

import (
	"stockz-app/config"
	"stockz-app/controller"
	"stockz-app/middlewares"

	"github.com/labstack/echo/v4"
	echotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/labstack/echo.v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func RouteInit() {
	tracer.Start(
		tracer.WithServiceName("stockz-service"),
		tracer.WithEnv("dev"),
	)
	defer tracer.Stop()
	e := echo.New()
	e.Use(echotrace.Middleware(echotrace.WithServiceName("stockz-service")))
	config.DatabaseInit()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello datadog")
	})
	app := e.Group("/api/v1")

	// User Authentication
	app.POST("/register", controller.Register)
	app.POST("/login", controller.Login)
	app.GET("/home", middlewares.Auth(controller.Home)) // Test authorization from jwt and cookies
	// Function for User
	app.POST("/stockz", middlewares.Auth(controller.CreatePost))
	app.GET("/stockz", controller.GetAllPost)
	app.GET("/stockz/:id", controller.PostById)
	app.PUT("/stockz/:id/edit", middlewares.Auth(controller.UpdatePost))
	app.POST("/stockz/:id/comment", middlewares.Auth(controller.CreateComment))
	app.GET("/stockz/profile/:username", controller.GetProfile)
	app.GET("/stockz/username", controller.SearchUsername)
	e.Logger.Fatal(e.Start(":8080"))
}
