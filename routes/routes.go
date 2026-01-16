package routes

import (
	"transfers/controllers"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// --- CONFIGURACIÓN GLOBAL ---
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false

	// Servir archivos estáticos
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
	// RUTAS PÚBLICAS
	// ---------------------------------------------------------
	r.GET("/", func(c *gin.Context) {
		c.File("./static/views/login.html")
	})
	r.POST("/login", authCtrl.Login)
	r.GET("/logout", authCtrl.Logout)

	// SEMÁFORO DE REDIRECCIÓN
	r.GET("/dashboard", utils.JWTAuthMiddleware(), authCtrl.RedirectByRole)

	// ---------------------------------------------------------
	// VISTAS DASHBOARD (HTML)
	// ---------------------------------------------------------
	vistas := r.Group("/")
	vistas.Use(utils.JWTAuthMiddleware())
	{
		// Panel Admin
		vistas.GET("/admin", func(c *gin.Context) {
			c.File("./static/views/admin.html")
		})

		// Módulo Usuarios
		usersGroup := vistas.Group("/dashboard/users")
		{
			uHandler := func(c *gin.Context) { c.File("./static/views/user.html") }
			usersGroup.GET("", uHandler)
			usersGroup.GET("/", uHandler)
			usersGroup.GET("/manage", func(c *gin.Context) {
				c.File("./static/views/user_crud.html")
			})
		}

		// Módulo Clientes
		clientsGroup := vistas.Group("/dashboard/clients")
		{
			clHandler := func(c *gin.Context) { c.File("./static/views/admin_client.html") }
			clientsGroup.GET("", clHandler)
			clientsGroup.GET("/", clHandler)
			clientsGroup.GET("/manage", func(c *gin.Context) {
				c.File("./static/views/admin_client_crud.html")
			})
		}

		// Otras vistas directas
		vistas.GET("/drivers", func(c *gin.Context) { c.File("./static/views/drivers.html") })
		vistas.GET("/companies", func(c *gin.Context) { c.File("./static/views/companies.html") })
	}

	// ---------------------------------------------------------
	// API V1 (Endpoints JSON)
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

		// Empresas
		api.GET("/companies", companyCtrl.GetAll)
		api.POST("/companies", companyCtrl.POST)
		api.GET("/companies/:id", companyCtrl.Get)
		api.PUT("/companies/:id", companyCtrl.PUT)
		api.DELETE("/companies/:id", companyCtrl.DELETE)

		// Conductores
		api.GET("/drivers", driverCtrl.GetAll)
		api.POST("/drivers", driverCtrl.POST)
		api.GET("/drivers/:id", driverCtrl.Get)
		api.PUT("/drivers/:id", driverCtrl.PUT)
		api.DELETE("/drivers/:id", driverCtrl.DELETE)

		// Vehículos
		api.GET("/vehicles", vehicleCtrl.GetAll)
		api.POST("/vehicles", vehicleCtrl.POST)
		api.GET("/vehicles/:id", vehicleCtrl.Get)
		api.PUT("/vehicles/:id", vehicleCtrl.PUT)
		api.DELETE("/vehicles/:id", vehicleCtrl.DELETE)

		// Reservas
		api.GET("/bookings", bookingCtrl.GetAll)
		api.POST("/bookings", bookingCtrl.POST)
		api.GET("/bookings/:id", bookingCtrl.Get)
		api.PUT("/bookings/:id", bookingCtrl.PUT)
		api.DELETE("/bookings/:id", bookingCtrl.DELETE)

		// Viajes (Rides) - USANDO rideCtrl
		api.GET("/rides", rideCtrl.GetAll)
		api.POST("/rides", rideCtrl.POST)
		api.GET("/rides/:id", rideCtrl.Get)
		api.PUT("/rides/:id", rideCtrl.PUT)
		api.DELETE("/rides/:id", rideCtrl.DELETE)

		// Pagos - USANDO paymentCtrl
		api.GET("/payments", paymentCtrl.GetAll)
		api.POST("/payments", paymentCtrl.POST)
		api.GET("/payments/:id", paymentCtrl.Get)

		// Eventos - USANDO eventCtrl
		api.GET("/events", eventCtrl.GetAll)
		api.POST("/events", eventCtrl.POST)

		// Valoraciones - USANDO ratingCtrl
		api.GET("/ratings", ratingCtrl.GetAll)
		api.POST("/ratings", ratingCtrl.POST)
	}

	return r
}
