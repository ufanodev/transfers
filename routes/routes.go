package routes

import (
	"fmt"
	"transfers/controllers"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// --- CONFIGURACIÓN GLOBAL ---
	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// --- SERVIR ESTÁTICOS ---
	// Mapeo general para que todas las subcarpetas de js, css e img sean accesibles
	r.Static("/js", "./static/js")
	r.Static("/css", "./static/css")
	r.Static("/img", "./static/img")

	// 1. Inicialización de controladores
	authCtrl := &controllers.AuthController{DB: db}
	userCtrl := &controllers.UserController{DB: db}
	clientCtrl := &controllers.ClientController{DB: db}
	companyCtrl := &controllers.CompanyController{DB: db}
	driverCtrl := &controllers.DriverController{DB: db}
	vehicleCtrl := &controllers.VehicleController{DB: db}
	bookingCtrl := &controllers.BookingController{DB: db}
	eventCtrl := &controllers.BookingEventController{DB: db}
	rideCtrl := &controllers.RideController{DB: db}
	paymentCtrl := &controllers.PaymentController{DB: db}
	ratingCtrl := &controllers.RatingController{DB: db}

	// ---------------------------------------------------------
	// RUTAS PÚBLICAS (Login / Logout)
	// ---------------------------------------------------------
	r.GET("/", func(c *gin.Context) {
		c.File("./static/views/auth/login.html")
	})
	r.POST("/login", authCtrl.Login)
	r.GET("/logout", authCtrl.Logout)

	// Dashboard inteligente (Semáforo de redirección)
	r.GET("/dashboard", utils.JWTAuthMiddleware(), authCtrl.RedirectByRole)

	// ---------------------------------------------------------
	// SECCIÓN ADMIN
	// ---------------------------------------------------------
	admin := r.Group("/admin")
	admin.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("admin"))
	{
		// Ruta raíz del grupo admin
		admin.GET("/", func(c *gin.Context) {
			fmt.Println("[ROUTER] 🖥️ Cargando Panel: Admin")
			c.File("./static/views/admin/admin.html")
		})
		admin.GET("/users", func(c *gin.Context) { c.File("./static/views/admin/admin_user.html") })
		admin.GET("/users/manage", func(c *gin.Context) { c.File("./static/views/admin/admin_user_crud.html") })
		admin.GET("/clients", func(c *gin.Context) { c.File("./static/views/admin/admin_client.html") })
		admin.GET("/clients/manage", func(c *gin.Context) { c.File("./static/views/admin/admin_client_crud.html") })
	}

	// ---------------------------------------------------------
	// SECCIÓN CLIENTE (Corregida para evitar 404)
	// ---------------------------------------------------------
	client := r.Group("/client")
	client.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("client"))
	{
		// Definimos la raíz del grupo claramente
		client.GET("/", func(c *gin.Context) {
			fmt.Println("[ROUTER] 👤 Cargando Panel: Cliente")
			c.Header("Content-Type", "text/html; charset=utf-8")
			c.File("./static/views/client/client.html")
		})

		// Sub-rutas de Cliente (Basadas en tu lista de opciones)
		client.GET("/profile", func(c *gin.Context) { c.File("./static/views/client/client_profile.html") })
		client.GET("/settings", func(c *gin.Context) { c.File("./static/views/client/client_settings.html") })
		client.GET("/account", func(c *gin.Context) { c.File("./static/views/client/client_account.html") })
		client.GET("/balance", func(c *gin.Context) { c.File("./static/views/client/client_balance.html") })
		client.GET("/payment-methods", func(c *gin.Context) { c.File("./static/views/client/client_payment_methods.html") })
		client.GET("/bookings", func(c *gin.Context) { c.File("./static/views/client/client_bookings.html") })
		client.GET("/bookings/new", func(c *gin.Context) { c.File("./static/views/client/client_booking_new.html") })
		client.GET("/rides", func(c *gin.Context) { c.File("./static/views/client/client_rides.html") })
		client.GET("/payments", func(c *gin.Context) { c.File("./static/views/client/client_payments.html") })
		client.GET("/invoices", func(c *gin.Context) { c.File("./static/views/client/client_invoices.html") })
		client.GET("/ratings", func(c *gin.Context) { c.File("./static/views/client/client_ratings.html") })
	}

	// ---------------------------------------------------------
	// SECCIÓN CONDUCTOR (Driver)
	// ---------------------------------------------------------
	driver := r.Group("/driver")
	driver.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("driver"))
	{
		driver.GET("/", func(c *gin.Context) {
			fmt.Println("[ROUTER] 🚗 Cargando Panel: Conductor")
			c.File("./static/views/driver/driver.html")
		})
	}

	// ---------------------------------------------------------
	// SECCIÓN EMPRESA (Company)
	// ---------------------------------------------------------
	company := r.Group("/company")
	company.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("company"))
	{
		company.GET("/", func(c *gin.Context) {
			fmt.Println("[ROUTER] 🏢 Cargando Panel: Empresa")
			c.File("./static/views/company/company.html")
		})
	}

	// ---------------------------------------------------------
	// API V1 (Endpoints de datos JSON)
	// ---------------------------------------------------------
	api := r.Group("/api/v1")
	api.Use(utils.JWTAuthMiddleware())
	{
		// Usuarios
		api.GET("/users", userCtrl.GetAll)
		api.POST("/users", userCtrl.POST)
		api.GET("/users/:id", userCtrl.Get)
		api.PUT("/users/:id", userCtrl.PUT)
		api.DELETE("/users/:id", userCtrl.DELETE)

		// Clientes
		api.GET("/clients", clientCtrl.GetAll)
		api.POST("/clients", clientCtrl.POST)
		api.GET("/clients/:id", clientCtrl.Get)
		api.PUT("/clients/:id", clientCtrl.PUT)
		api.DELETE("/clients/:id", clientCtrl.DELETE)

		// Empresas y Conductores
		api.GET("/companies", companyCtrl.GetAll)
		api.POST("/companies", companyCtrl.POST)
		api.GET("/drivers", driverCtrl.GetAll)
		api.POST("/drivers", driverCtrl.POST)

		// Flota y Reservas
		api.GET("/vehicles", vehicleCtrl.GetAll)
		api.POST("/vehicles", vehicleCtrl.POST)
		api.GET("/bookings", bookingCtrl.GetAll)
		api.POST("/bookings", bookingCtrl.POST)
		api.GET("/rides", rideCtrl.GetAll)
		api.POST("/rides", rideCtrl.POST)

		// Otros
		api.GET("/payments", paymentCtrl.GetAll)
		api.GET("/events", eventCtrl.GetAll)
		api.GET("/ratings", ratingCtrl.GetAll)
	}

	return r
}
