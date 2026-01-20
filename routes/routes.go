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
	r.Static("/js", "static/js")
	r.Static("/css", "static/css")
	r.Static("/img", "static/img")

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
		c.File("static/views/auth/login.html")
	})
	r.POST("/login", authCtrl.Login)
	r.GET("/logout", authCtrl.Logout)

	// Semáforo de redirección inicial tras login
	r.GET("/dashboard", utils.JWTAuthMiddleware(), authCtrl.RedirectByRole)

	// ---------------------------------------------------------
	// SECCIÓN ADMIN (Prefijo /admin/)
	// ---------------------------------------------------------
	admin := r.Group("/admin")
	admin.Use(utils.JWTAuthMiddleware(), utils.RoleMiddleware("admin"))
	{
		// 🏠 Inicio Admin
		admin.GET("/", func(c *gin.Context) {
			c.File("static/views/admin/admin.html")
		})

		// 👥 Gestión de Usuarios
		admin.GET("/users", func(c *gin.Context) {
			fmt.Println("[ADMIN] -> Cargando Vista Usuarios")
			c.File("static/views/admin/admin_user.html")
		})
		admin.GET("/users/manage", func(c *gin.Context) {
			c.File("static/views/admin/admin_user_crud.html")
		})

		// 💼 Gestión de Clientes
		admin.GET("/clients", func(c *gin.Context) {
			fmt.Println("[ADMIN] -> Cargando Vista Clientes")
			c.File("static/views/admin/admin_client.html")
		})
		admin.GET("/clients/manage", func(c *gin.Context) {
			c.File("static/views/admin/admin_client_crud.html")
		})

		// 🚗 Rutas Preparadas para futuras vistas
		admin.GET("/drivers", func(c *gin.Context) { c.String(200, "Vista de Conductores") })
		admin.GET("/vehicles", func(c *gin.Context) { c.String(200, "Vista de Vehículos") })
		admin.GET("/companies", func(c *gin.Context) { c.String(200, "Vista de Empresas") })
		admin.GET("/bookings", func(c *gin.Context) { c.String(200, "Vista de Reservas") })
		admin.GET("/rides", func(c *gin.Context) { c.String(200, "Vista de Viajes") })
		admin.GET("/payments", func(c *gin.Context) { c.String(200, "Vista de Pagos") })
		admin.GET("/ratings", func(c *gin.Context) { c.String(200, "Vista de Valoraciones") })
		admin.GET("/events", func(c *gin.Context) { c.String(200, "Vista de Eventos") })
	}

	// ---------------------------------------------------------
	// SECCIONES OTROS ROLES
	// ---------------------------------------------------------
	r.GET("/client/", utils.JWTAuthMiddleware(), utils.RoleMiddleware("client"), func(c *gin.Context) {
		c.File("static/views/client/client.html")
	})
	r.GET("/driver/", utils.JWTAuthMiddleware(), utils.RoleMiddleware("driver"), func(c *gin.Context) {
		c.File("static/views/driver/driver.html")
	})
	r.GET("/company/", utils.JWTAuthMiddleware(), utils.RoleMiddleware("company"), func(c *gin.Context) {
		c.File("static/views/company/company.html")
	})

	// ---------------------------------------------------------
	// API V1 (Endpoints JSON)
	// ---------------------------------------------------------
	api := r.Group("/api/v1")
	api.Use(utils.JWTAuthMiddleware())
	{
		api.GET("/users/me", userCtrl.GetMe)

		// 👥 CRUD USUARIOS
		api.GET("/users", userCtrl.GetAll)
		api.POST("/users", userCtrl.POST)
		api.GET("/users/:id", userCtrl.Get)
		api.PUT("/users/:id", userCtrl.PUT)
		api.DELETE("/users/:id", userCtrl.DELETE)

		// 💼 CRUD CLIENTES (Actualizado con ID y métodos completos)
		api.GET("/clients", clientCtrl.GetAll)
		api.POST("/clients", clientCtrl.POST)         // Crear Cliente
		api.GET("/clients/:id", clientCtrl.Get)       // Carga datos en CRUD
		api.PUT("/clients/:id", clientCtrl.PUT)       // Actualizar Cliente
		api.DELETE("/clients/:id", clientCtrl.DELETE) // Eliminar Cliente

		// 🚗 Endpoints de Datos Generales
		api.GET("/drivers", driverCtrl.GetAll)
		api.GET("/vehicles", vehicleCtrl.GetAll)
		api.GET("/companies", companyCtrl.GetAll)
		api.GET("/bookings", bookingCtrl.GetAll)
		api.GET("/rides", rideCtrl.GetAll)
		api.GET("/payments", paymentCtrl.GetAll)
		api.GET("/events", eventCtrl.GetAll)
		api.GET("/ratings", ratingCtrl.GetAll)
	}

	return r
}
