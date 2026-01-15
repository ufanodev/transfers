package routes

import (
	"net/http"
	"transfers/controllers"
	"transfers/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// Servir archivos estáticos (JS y CSS)
	r.Static("/js", "./static/js")
	r.Static("/css", "./static/css")

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
	r.GET("/", func(c *gin.Context) { c.File("./static/views/login.html") })
	r.POST("/login", authCtrl.Login)
	r.GET("/logout", func(c *gin.Context) {
		utils.ClearAuthCookie(c)
		c.Redirect(http.StatusSeeOther, "/")
	})

	// ---------------------------------------------------------
	// VISTAS DASHBOARD (Protegidas por JWT y Rol Admin)
	// ---------------------------------------------------------
	dashboard := r.Group("/dashboard")
	dashboard.Use(utils.JWTAuthMiddleware())
	{
		// Redirección principal automática según rol
		dashboard.GET("/", authCtrl.RedirectByRole)

		// --- MÓDULO USUARIOS ---
		users := dashboard.Group("/users")
		{
			// Lista principal de usuarios
			users.GET("/", func(c *gin.Context) {
				if role, _ := c.Get("userRole"); role != "admin" {
					c.Redirect(http.StatusSeeOther, "/dashboard/")
					return
				}
				c.File("./static/views/user.html")
			})

			// Formulario Crear/Editar/Borrar (user_crud.html)
			users.GET("/manage", func(c *gin.Context) {
				if role, _ := c.Get("userRole"); role != "admin" {
					c.Redirect(http.StatusSeeOther, "/dashboard/")
					return
				}
				c.File("./static/views/user_crud.html")
			})
		}

		// Aquí puedes ir añadiendo los grupos de vistas para /clients, /companies, etc.
	}

	// ---------------------------------------------------------
	// API V1 (Endpoints JSON)
	// ---------------------------------------------------------
	api := r.Group("/api/v1")
	api.Use(utils.JWTAuthMiddleware())
	{
		// Usuarios
		u := api.Group("/users")
		{
			u.GET("/", userCtrl.GetAll)
			u.POST("/", userCtrl.POST)
			u.GET("/:id", userCtrl.Get)
			u.PUT("/:id", userCtrl.PUT)
			u.DELETE("/:id", userCtrl.DELETE)
		}

		// Clientes
		cl := api.Group("/clients")
		{
			cl.GET("/", clientCtrl.GetAll)
			cl.POST("/", clientCtrl.POST)
			cl.GET("/:id", clientCtrl.Get)
			cl.PUT("/:id", clientCtrl.PUT)
			cl.DELETE("/:id", clientCtrl.DELETE)
		}

		// Empresas
		co := api.Group("/companies")
		{
			co.GET("/", companyCtrl.GetAll)
			co.POST("/", companyCtrl.POST)
			co.GET("/:id", companyCtrl.Get)
			co.PUT("/:id", companyCtrl.PUT)
			co.DELETE("/:id", companyCtrl.DELETE)
		}

		// Conductores
		dr := api.Group("/drivers")
		{
			dr.GET("/", driverCtrl.GetAll)
			dr.POST("/", driverCtrl.POST)
			dr.GET("/:id", driverCtrl.Get)
			dr.PUT("/:id", driverCtrl.PUT)
			dr.DELETE("/:id", driverCtrl.DELETE)
		}

		// Vehículos
		vh := api.Group("/vehicles")
		{
			vh.GET("/", vehicleCtrl.GetAll)
			vh.POST("/", vehicleCtrl.POST)
			vh.GET("/:id", vehicleCtrl.Get)
			vh.PUT("/:id", vehicleCtrl.PUT)
			vh.DELETE("/:id", vehicleCtrl.DELETE)
		}

		// Reservas
		bk := api.Group("/bookings")
		{
			bk.GET("/", bookingCtrl.GetAll)
			bk.POST("/", bookingCtrl.POST)
			bk.GET("/:id", bookingCtrl.Get)
			bk.PUT("/:id", bookingCtrl.PUT)
			bk.DELETE("/:id", bookingCtrl.DELETE)
		}

		// Eventos de Reservas
		ev := api.Group("/events")
		{
			ev.GET("/", eventCtrl.GetAll)
			ev.POST("/", eventCtrl.POST)
			ev.GET("/:id", eventCtrl.Get)
			ev.PUT("/:id", eventCtrl.PUT)
			ev.DELETE("/:id", eventCtrl.DELETE)
		}

		// Viajes (Rides)
		rd := api.Group("/rides")
		{
			rd.GET("/", rideCtrl.GetAll)
			rd.POST("/", rideCtrl.POST)
			rd.GET("/:id", rideCtrl.Get)
			rd.PUT("/:id", rideCtrl.PUT)
			rd.DELETE("/:id", rideCtrl.DELETE)
		}

		// Pagos
		py := api.Group("/payments")
		{
			py.GET("/", paymentCtrl.GetAll)
			py.POST("/", paymentCtrl.POST)
			py.GET("/:id", paymentCtrl.Get)
			py.PUT("/:id", paymentCtrl.PUT)
			py.DELETE("/:id", paymentCtrl.DELETE)
		}

		// Valoraciones
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
