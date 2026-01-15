package routes

import (
	"transfers/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 1. Inicialización de todos los controladores
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

	// 2. Definición del grupo base de la API
	api := r.Group("/api/v1")
	{
		// Endpoints de Usuarios
		users := api.Group("/users")
		{
			users.GET("/", userCtrl.GetAll)
			users.GET("/:id", userCtrl.Get)
			users.POST("/", userCtrl.POST)
			users.PUT("/:id", userCtrl.PUT)
			users.DELETE("/:id", userCtrl.DELETE)
		}

		// Endpoints de Clientes
		clients := api.Group("/clients")
		{
			clients.GET("/", clientCtrl.GetAll)
			clients.GET("/:id", clientCtrl.Get)
			clients.POST("/", clientCtrl.POST)
			clients.PUT("/:id", clientCtrl.PUT)
			clients.DELETE("/:id", clientCtrl.DELETE)
		}

		// Endpoints de Empresas
		companies := api.Group("/companies")
		{
			companies.GET("/", companyCtrl.GetAll)
			companies.GET("/:id", companyCtrl.Get)
			companies.POST("/", companyCtrl.POST)
			companies.PUT("/:id", companyCtrl.PUT)
			companies.DELETE("/:id", companyCtrl.DELETE)
		}

		// Endpoints de Conductores
		drivers := api.Group("/drivers")
		{
			drivers.GET("/", driverCtrl.GetAll)
			drivers.GET("/:id", driverCtrl.Get)
			drivers.POST("/", driverCtrl.POST)
			drivers.PUT("/:id", driverCtrl.PUT)
			drivers.DELETE("/:id", driverCtrl.DELETE)
		}

		// Endpoints de Vehículos
		vehicles := api.Group("/vehicles")
		{
			vehicles.GET("/", vehicleCtrl.GetAll)
			vehicles.GET("/:id", vehicleCtrl.Get)
			vehicles.POST("/", vehicleCtrl.POST)
			vehicles.PUT("/:id", vehicleCtrl.PUT)
			vehicles.DELETE("/:id", vehicleCtrl.DELETE)
		}

		// Endpoints de Reservas (Bookings)
		bookings := api.Group("/bookings")
		{
			bookings.GET("/", bookingCtrl.GetAll)
			bookings.GET("/:id", bookingCtrl.Get)
			bookings.POST("/", bookingCtrl.POST)
			bookings.PUT("/:id", bookingCtrl.PUT)
			bookings.DELETE("/:id", bookingCtrl.DELETE)
		}

		// Endpoints de Eventos de Reserva
		events := api.Group("/events")
		{
			events.GET("/", eventCtrl.GetAll)
			events.GET("/:id", eventCtrl.Get)
			events.POST("/", eventCtrl.POST)
			events.PUT("/:id", eventCtrl.PUT)
			events.DELETE("/:id", eventCtrl.DELETE)
		}

		// Endpoints de Viajes (Rides)
		rides := api.Group("/rides")
		{
			rides.GET("/", rideCtrl.GetAll)
			rides.GET("/:id", rideCtrl.Get)
			rides.POST("/", rideCtrl.POST)
			rides.PUT("/:id", rideCtrl.PUT)
			rides.DELETE("/:id", rideCtrl.DELETE)
		}

		// Endpoints de Pagos
		payments := api.Group("/payments")
		{
			payments.GET("/", paymentCtrl.GetAll)
			payments.GET("/:id", paymentCtrl.Get)
			payments.POST("/", paymentCtrl.POST)
			payments.PUT("/:id", paymentCtrl.PUT)
			payments.DELETE("/:id", paymentCtrl.DELETE)
		}

		// Endpoints de Valoraciones (Ratings)
		ratings := api.Group("/ratings")
		{
			ratings.GET("/", ratingCtrl.GetAll)
			ratings.GET("/:id", ratingCtrl.Get)
			ratings.POST("/", ratingCtrl.POST)
			ratings.PUT("/:id", ratingCtrl.PUT)
			ratings.DELETE("/:id", ratingCtrl.DELETE)
		}
	}

	return r
}
