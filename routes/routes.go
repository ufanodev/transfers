package routes

import (
	"transfers/controllers"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.RedirectTrailingSlash = true
	r.RedirectFixedPath = true

	// --- SERVIR ARCHIVOS ESTÁTICOS ---
	r.Static("/js", "static/js")
	r.Static("/css", "static/css")
	r.Static("/img", "static/img")

	// 1. INICIALIZACIÓN DE CONTROLADORES
	authCtrl := &controllers.AuthController{DB: db}
	userCtrl := &controllers.UserController{DB: db}
	clientCtrl := &controllers.ClientController{DB: db}
	companyCtrl := &controllers.CompanyController{DB: db}
	driverCtrl := &controllers.DriverController{DB: db}
	vehicleCtrl := &controllers.VehicleController{DB: db}
	rideCtrl := &controllers.RideController{DB: db}
	paymentCtrl := &controllers.PaymentController{DB: db}
	bookingCtrl := &controllers.BookingController{DB: db} // Nuevo controlador

	// ---------------------------------------------------------
	// RUTAS PÚBLICAS
	// ---------------------------------------------------------
	r.GET("/", func(c *gin.Context) { c.File("static/views/auth/login.html") })
	r.POST("/login", authCtrl.Login)
	r.GET("/logout", authCtrl.Logout)
	r.GET("/dashboard", utils.JWTAuthMiddleware(), authCtrl.RedirectByRole)

	// ---------------------------------------------------------
	// SECCIÓN ADMINISTRACIÓN (Vistas HTML)
	// ---------------------------------------------------------
	admin := r.Group("/admin")
	admin.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("admin"))
	{
		admin.GET("/", func(c *gin.Context) { c.File("static/views/admin/admin.html") })

		// Usuarios
		admin.GET("/users", func(c *gin.Context) { c.File("static/views/admin/admin_user.html") })
		admin.GET("/users/manage", func(c *gin.Context) { c.File("static/views/admin/admin_user_crud.html") })

		// Clientes
		admin.GET("/clients", func(c *gin.Context) { c.File("static/views/admin/admin_client.html") })
		admin.GET("/clients/manage", func(c *gin.Context) { c.File("static/views/admin/admin_client_crud.html") })

		// Empresas
		admin.GET("/companies", func(c *gin.Context) { c.File("static/views/admin/admin_company.html") })
		admin.GET("/companies/manage", func(c *gin.Context) { c.File("static/views/admin/admin_company_crud.html") })

		// Conductores (Drivers)
		admin.GET("/drivers", func(c *gin.Context) { c.File("static/views/admin/admin_driver.html") })
		admin.GET("/drivers/manage", func(c *gin.Context) { c.File("static/views/admin/admin_driver_crud.html") })

		// Vehículos
		admin.GET("/vehicles", func(c *gin.Context) { c.File("static/views/admin/admin_vehicle.html") })
		admin.GET("/vehicles/manage", func(c *gin.Context) { c.File("static/views/admin/admin_vehicle_crud.html") })

		// Reservas y Logística
		admin.GET("/bookings", func(c *gin.Context) { c.File("static/views/admin/admin_booking.html") })
		admin.GET("/rides", func(c *gin.Context) { c.File("static/views/admin/admin_ride.html") })
		admin.GET("/payments", func(c *gin.Context) { c.File("static/views/admin/admin_payment.html") })
	}

	// ---------------------------------------------------------
	// API V1 (Endpoints JSON)
	// ---------------------------------------------------------
	api := r.Group("/api/v1")
	api.Use(utils.JWTAuthMiddleware())
	{
		// CRUD Usuarios
		api.GET("/users", userCtrl.GetAll)
		api.POST("/users", userCtrl.POST)
		api.GET("/users/:id", userCtrl.Get)
		api.PUT("/users/:id", userCtrl.PUT)
		api.DELETE("/users/:id", userCtrl.DELETE)

		// CRUD Clientes
		api.GET("/clients", clientCtrl.GetAll)
		api.POST("/clients", clientCtrl.POST)
		api.GET("/clients/:id", clientCtrl.Get)
		api.PUT("/clients/:id", clientCtrl.PUT)
		api.DELETE("/clients/:id", clientCtrl.DELETE)

		// CRUD Empresas
		api.GET("/companies", companyCtrl.GetAll)
		api.POST("/companies", companyCtrl.POST)
		api.GET("/companies/:id", companyCtrl.Get)
		api.PUT("/companies/:id", companyCtrl.PUT)
		api.DELETE("/companies/:id", companyCtrl.DELETE)

		// CRUD Conductores
		api.GET("/drivers", driverCtrl.GetAll)
		api.POST("/drivers", driverCtrl.POST)
		api.GET("/drivers/:id", driverCtrl.Get)
		api.PUT("/drivers/:id", driverCtrl.PUT)
		api.DELETE("/drivers/:id", driverCtrl.DELETE)

		// CRUD Vehículos
		api.GET("/vehicles", vehicleCtrl.GetAll)
		api.POST("/vehicles", vehicleCtrl.POST)
		api.GET("/vehicles/:id", vehicleCtrl.Get)
		api.PUT("/vehicles/:id", vehicleCtrl.PUT)
		api.DELETE("/vehicles/:id", vehicleCtrl.DELETE)

		// Gestión de Viajes y Pagos
		api.GET("/rides", rideCtrl.GetAll)
		api.POST("/rides", rideCtrl.POST)
		api.PUT("/rides/:id", rideCtrl.PUT)

		api.GET("/payments", paymentCtrl.GetAll)
		api.POST("/payments", paymentCtrl.POST)

		api.GET("/bookings", bookingCtrl.GetAll)
		api.POST("/bookings", bookingCtrl.POST)
		api.GET("/bookings/:id", bookingCtrl.Get)
		api.PUT("/bookings/:id", bookingCtrl.PUT)
	}

	return r
}
