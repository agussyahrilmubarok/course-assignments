package server

import (
	"finly-go/internal/controller"

	"github.com/gofiber/fiber/v2"
)

func WebRoutes(app *fiber.App) {

	collections := GetCollections()

	homeController := controller.NewHomeController()
	authController := controller.NewAuthController(collections.UserCollection)
	dashboardController := controller.NewDashboardController(collections.InvoiceCollection, collections.CustomerCollection)
	customerController := controller.NewCustomerController(collections.CustomerCollection)
	invoiceController := controller.NewInvoiceController(collections.InvoiceCollection, collections.CustomerCollection)

	app.Get("/", homeController.GetIndex)
	app.Get("/signup", PreventAfterAuth, authController.GetSignUp)
	app.Post("/signup", PreventAfterAuth, authController.PostSignUp)
	app.Get("/login", PreventAfterAuth, authController.GetLogin)
	app.Post("/login", PreventAfterAuth, authController.PostLogin)
	app.Get("/logout", RequireAuth, authController.GetLogout)

	app.Get("/dashboard", RequireAuth, dashboardController.GetDashboard)
	app.Get("/dashboard/customers", RequireAuth, customerController.Index)
	app.Get("/dashboard/customers/new", RequireAuth, customerController.GetNew)
	app.Post("/dashboard/customers/new", RequireAuth, customerController.PostNew)
	app.Get("/dashboard/customers/:id/edit", RequireAuth, customerController.GetEdit)
	app.Post("/dashboard/customers/:id/edit", RequireAuth, customerController.PostEdit)
	app.Post("/dashboard/customers/:id/delete", RequireAuth, customerController.PostDelete)

	app.Get("/dashboard/invoices", RequireAuth, invoiceController.Index)
	app.Get("/dashboard/invoices/new", RequireAuth, invoiceController.GetNew)
	app.Post("/dashboard/invoices/new", RequireAuth, invoiceController.PostNew)
	app.Get("/dashboard/invoices/:id/edit", RequireAuth, invoiceController.GetEdit)
	app.Post("/dashboard/invoices/:id/edit", RequireAuth, invoiceController.PostEdit)
	app.Post("/dashboard/invoices/:id/delete", RequireAuth, invoiceController.PostDelete)
}
