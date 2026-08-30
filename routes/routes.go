package routes

import (
	"bazar-be/controllers"
	"bazar-be/dto"
	"bazar-be/middleware"

	"github.com/gofiber/fiber/v3"
)

// SetupRoutes registers all application routes
func SetupRoutes(app *fiber.App) {
	authController := controllers.NewAuthController()
	categoryController := controllers.NewCategoryController()
	vendorController := controllers.NewVendorController()
	vendorOutletController := controllers.NewVendorOutletController()
	bazaarController := controllers.NewBazaarController()
	outletController := controllers.NewOutletController()
	productController := controllers.NewProductController()
	transactionController := controllers.NewTransactionController()
	transactionItemController := controllers.NewTransactionItemController()
	devController := controllers.NewDevController()

	// Base / Health check route
	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(dto.APIResponse{
			Success: true,
			Message: "Bazar Backend API is running",
		})
	})

	api := app.Group("/api")

	// ==========================================
	// Auth Routes
	// ==========================================
	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/logout", middleware.AuthMiddleware(), authController.Logout)
	auth.Get("/me", middleware.AuthMiddleware(), authController.GetProfile)

	// ==========================================
	// Category Routes
	// ==========================================
	categories := api.Group("/categories")
	categories.Get("/", categoryController.GetAll)
	categories.Get("/:id", categoryController.GetByID)
	categories.Post("/", middleware.AuthMiddleware(), categoryController.Create)
	categories.Put("/:id", middleware.AuthMiddleware(), categoryController.Update)
	categories.Delete("/:id", middleware.AuthMiddleware(), categoryController.Delete)

	// ==========================================
	// Vendor Routes
	// ==========================================
	vendors := api.Group("/vendors")
	vendors.Get("/", vendorController.GetAll)
	vendors.Get("/:id", vendorController.GetByID)
	vendors.Post("/", middleware.AuthMiddleware(), vendorController.Create)
	vendors.Put("/:id", middleware.AuthMiddleware(), vendorController.Update)
	vendors.Delete("/:id", middleware.AuthMiddleware(), vendorController.Delete)

	// ==========================================
	// Vendor-Outlet Routes
	// ==========================================
	vendorOutlets := api.Group("/vendor-outlets")
	vendorOutlets.Get("/", vendorOutletController.GetAll)
	vendorOutlets.Get("/:id", vendorOutletController.GetByID)
	vendorOutlets.Post("/", middleware.AuthMiddleware(), vendorOutletController.Assign)
	vendorOutlets.Put("/:id", middleware.AuthMiddleware(), vendorOutletController.Update)
	vendorOutlets.Delete("/:id", middleware.AuthMiddleware(), vendorOutletController.Delete)

	// ==========================================
	// Bazaar Routes
	// ==========================================
	bazaars := api.Group("/bazaars")
	bazaars.Get("/", bazaarController.GetAll)
	bazaars.Get("/:id", bazaarController.GetByID)
	bazaars.Post("/", middleware.AuthMiddleware(), bazaarController.Create)
	bazaars.Put("/:id", middleware.AuthMiddleware(), bazaarController.Update)
	bazaars.Delete("/:id", middleware.AuthMiddleware(), bazaarController.Delete)

	// ==========================================
	// Outlet Routes
	// ==========================================
	outlets := api.Group("/outlets")
	outlets.Get("/", outletController.GetAll)
	outlets.Get("/:id", outletController.GetByID)
	outlets.Post("/", middleware.AuthMiddleware(), outletController.Create)
	outlets.Put("/:id", middleware.AuthMiddleware(), outletController.Update)
	outlets.Delete("/:id", middleware.AuthMiddleware(), outletController.Delete)

	// ==========================================
	// Product Routes
	// ==========================================
	products := api.Group("/products")
	products.Get("/", productController.GetAll)
	products.Get("/:id", productController.GetByID)
	products.Post("/", middleware.AuthMiddleware(), productController.Create)
	products.Put("/:id", middleware.AuthMiddleware(), productController.Update)
	products.Delete("/:id", middleware.AuthMiddleware(), productController.Delete)

	// ==========================================
	// Transaction Routes
	// ==========================================
	transactions := api.Group("/transactions")
	transactions.Get("/summary", middleware.AuthMiddleware(), transactionController.GetSummary)
	transactions.Get("/", middleware.AuthMiddleware(), transactionController.GetAll)
	transactions.Get("/:id", middleware.AuthMiddleware(), transactionController.GetByID)
	transactions.Post("/", middleware.AuthMiddleware(), transactionController.Create)

	// ==========================================
	// Transaction Items (Sales Analytics / Detail)
	// ==========================================
	transactionItems := api.Group("/transaction-items")
	transactionItems.Get("/top-selling", middleware.AuthMiddleware(), transactionItemController.GetTopSelling)
	transactionItems.Get("/", middleware.AuthMiddleware(), transactionItemController.GetAll)
	transactionItems.Get("/:id", middleware.AuthMiddleware(), transactionItemController.GetByID)

	// ==========================================
	// Dev & Database Seeder Routes
	// ==========================================
	dev := api.Group("/dev")
	dev.Post("/seed", devController.Seed)
	dev.Get("/seed", devController.Seed)
}
