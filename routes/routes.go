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

	// Servir archivos estáticos (JS, CSS, Imágenes)
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
		fmt.Println("[SERVER LOG] Sirviendo vista: login.html")
		c.File("./static/views/login.html")
	})
	r.POST("/login", authCtrl.Login)
	r.GET("/logout", authCtrl.Logout)

	// ---------------------------------------------------------
	// SEMÁFORO DE REDIRECCIÓN (Tras Login exitoso)
	// ---------------------------------------------------------
	// El JS de login debe redirigir a /dashboard para que este decida el destino
	r.GET("/dashboard", utils.JWTAuthMiddleware(), authCtrl.RedirectByRole)

	// ---------------------------------------------------------
	// VISTAS RAÍZ Y DASHBOARD (Protegidas por JWT)
	// ---------------------------------------------------------
	vistas := r.Group("/")
	vistas.Use(utils.JWTAuthMiddleware())
	{
		// --- DASHBOARDS PRINCIPALES ---

		// Admin: http://localhost:8080/admin
		vistas.GET("/admin", func(c *gin.Context) {
			fmt.Println("[VIEW LOG] Entregando: admin.html")
			c.File("./static/views/admin.html")
		})

		// Drivers: http://localhost:8080/drivers
		vistas.GET("/drivers", func(c *gin.Context) {
			fmt.Println("[VIEW LOG] Entregando: drivers.html")
			c.File("./static/views/drivers.html")
		})

		// Clients: http://localhost:8080/clients
		vistas.GET("/clients", func(c *gin.Context) {
			fmt.Println("[VIEW LOG] Entregando: client.html")
			c.File("./static/views/client.html")
		})

		// Companies: http://localhost:8080/companies
		vistas.GET("/companies", func(c *gin.Context) {
			fmt.Println("[VIEW LOG] Entregando: companies.html")
			c.File("./static/views/companies.html")
		})

		// --- MÓDULO DE USUARIOS (Soporta /dashboard/users solicitado por admin.html) ---

		usersGroup := vistas.Group("/dashboard/users")
		{
			// Lista: http://localhost:8080/dashboard/users
			usersGroup.GET("/", func(c *gin.Context) {
				fmt.Println("[VIEW LOG] Entregando: user.html (Lista)")
				c.File("./static/views/user.html")
			})

			// CRUD: http://localhost:8080/dashboard/users/manage
			usersGroup.GET("/manage", func(c *gin.Context) {
				fmt.Println("[VIEW LOG] Entregando: user_crud.html (Formulario)")
				c.File("./static/views/user_crud.html")
			})
		}
	}

	// ---------------------------------------------------------
	// API V1 (Endpoints de Datos JSON)
	// ---------------------------------------------------------
	api := r.Group("/api/v1")
	api.Use(utils.JWTAuthMiddleware())
	{
		// API de Usuarios
		u := api.Group("/users")
		{
			u.GET("/", userCtrl.GetAll)
			u.POST("/", userCtrl.POST)
			u.GET("/:id", userCtrl.Get)
			u.PUT("/:id", userCtrl.PUT)
			u.DELETE("/:id", userCtrl.DELETE)
		}

		// API de Clientes
		cl := api.Group("/clients")
		{
			cl.GET("/", clientCtrl.GetAll)
			cl.POST("/", clientCtrl.POST)
			cl.GET("/:id", clientCtrl.Get)
			cl.PUT("/:id", clientCtrl.PUT)
			cl.DELETE("/:id", clientCtrl.DELETE)
		}

		// API de Empresas
		co := api.Group("/companies")
		{
			co.GET("/", companyCtrl.GetAll)
			co.POST("/", companyCtrl.POST)
			co.GET("/:id", companyCtrl.Get)
			co.PUT("/:id", companyCtrl.PUT)
			co.DELETE("/:id", companyCtrl.DELETE)
		}

		// API de Conductores
		dr := api.Group("/drivers")
		{
			dr.GET("/", driverCtrl.GetAll)
			dr.POST("/", driverCtrl.POST)
			dr.GET("/:id", driverCtrl.Get)
			dr.PUT("/:id", driverCtrl.PUT)
			dr.DELETE("/:id", driverCtrl.DELETE)
		}

		// API de Vehículos
		vh := api.Group("/vehicles")
		{
			vh.GET("/", vehicleCtrl.GetAll)
			vh.POST("/", vehicleCtrl.POST)
			vh.GET("/:id", vehicleCtrl.Get)
			vh.PUT("/:id", vehicleCtrl.PUT)
			vh.DELETE("/:id", vehicleCtrl.DELETE)
		}

		// API de Reservas
		bk := api.Group("/bookings")
		{
			bk.GET("/", bookingCtrl.GetAll)
			bk.POST("/", bookingCtrl.POST)
			bk.GET("/:id", bookingCtrl.Get)
			bk.PUT("/:id", bookingCtrl.PUT)
			bk.DELETE("/:id", bookingCtrl.DELETE)
		}

		// API de Eventos
		ev := api.Group("/events")
		{
			ev.GET("/", eventCtrl.GetAll)
			ev.POST("/", eventCtrl.POST)
			ev.GET("/:id", eventCtrl.Get)
			ev.PUT("/:id", eventCtrl.PUT)
			ev.DELETE("/:id", eventCtrl.DELETE)
		}

		// API de Viajes (Rides)
		rd := api.Group("/rides")
		{
			rd.GET("/", rideCtrl.GetAll)
			rd.POST("/", rideCtrl.POST)
			rd.GET("/:id", rideCtrl.Get)
			rd.PUT("/:id", rideCtrl.PUT)
			rd.DELETE("/:id", rideCtrl.DELETE)
		}

		// API de Pagos
		py := api.Group("/payments")
		{
			py.GET("/", paymentCtrl.GetAll)
			py.POST("/", paymentCtrl.POST)
			py.GET("/:id", paymentCtrl.Get)
			py.PUT("/:id", paymentCtrl.PUT)
			py.DELETE("/:id", paymentCtrl.DELETE)
		}

		// API de Valoraciones
		rt := api.Group("/ratings")
		{
			rt.GET("/", ratingCtrl.GetAll)
			rt.POST("/", ratingCtrl.POST)
			rt.GET("/:id", ratingCtrl.Get)
			rt.PUT("/:id", ratingCtrl.PUT)
			rt.DELETE("/:id", ratingCtrl.DELETE)
		}
	}

	return r
}
