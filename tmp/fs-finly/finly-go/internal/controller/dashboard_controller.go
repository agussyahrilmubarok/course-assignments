package controller

import (
	"context"
	"encoding/json"
	"finly-go/internal/model"
	"finly-go/internal/sessions"
	"finly-go/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/leekchan/accounting"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type DashboardController struct {
	invoiceCollection  *mongo.Collection
	customerCollection *mongo.Collection
}

func NewDashboardController(invoiceCollection, customerCollection *mongo.Collection) *DashboardController {
	return &DashboardController{
		invoiceCollection:  invoiceCollection,
		customerCollection: customerCollection,
	}
}

func (dc *DashboardController) GetDashboard(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Dashboard",
	}

	// Get the user profile ID from session
	profileID, err := sessions.GetProfile(c)
	if err != nil || profileID == "" {
		logger.Log.Error("Error getting profile from session", zap.Error(err))
		return c.Redirect("/login")
	}
	data["ProfileID"] = profileID

	// Prepare query to fetch invoices and customers for the logged-in user
	query := map[string]interface{}{"owner": profileID}

	// Fetch the number of invoices
	invoiceCount, err := dc.invoiceCollection.CountDocuments(context.TODO(), query)
	if err != nil {
		logger.Log.Error("Error fetching invoice count", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching invoice count")
	}

	// Fetch the number of customers
	customerCount, err := dc.customerCollection.CountDocuments(context.TODO(), query)
	if err != nil {
		logger.Log.Error("Error fetching customer count", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching customer count")
	}

	// Fetch all invoices for the user
	cursor, err := dc.invoiceCollection.Find(context.TODO(), query, options.Find().SetSort(map[string]interface{}{"date": -1}))
	if err != nil {
		logger.Log.Error("Error fetching invoices", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching invoices")
	}
	defer cursor.Close(context.TODO())

	// Variables to calculate total paid and total pending
	var totalPaid, totalPending float64
	var allInvoices []model.Invoice

	// Iterate through invoices to calculate totals and build invoice list
	for cursor.Next(context.TODO()) {
		var invoice model.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			logger.Log.Error("Error decoding invoice", zap.Error(err))
			continue
		}
		allInvoices = append(allInvoices, invoice)
		if invoice.Status == "paid" {
			totalPaid += invoice.Amount
		} else if invoice.Status == "pending" {
			totalPending += invoice.Amount
		}
	}

	// Fetch the latest 5 invoices
	var latestInvoices []model.Invoice
	if len(allInvoices) > 5 {
		latestInvoices = allInvoices[:5]
	} else {
		latestInvoices = allInvoices
	}

	// Fetch customer data
	customerMap := make(map[string]string)
	cursorCustomers, err := dc.customerCollection.Find(context.TODO(), query)
	if err != nil {
		logger.Log.Error("Error fetching customers", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching customers")
	}
	defer cursorCustomers.Close(context.TODO())

	var customers []model.Customer
	if err := cursorCustomers.All(context.TODO(), &customers); err != nil {
		logger.Log.Error("Error decoding customer data", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error decoding customer data")
	}

	// Map customer IDs to their names
	for _, customer := range customers {
		customerMap[customer.ID] = customer.Name
	}

	// Associate CustomerName with each invoice
	for i := range allInvoices {
		if name, ok := customerMap[allInvoices[i].Customer]; ok {
			allInvoices[i].CustomerName = name
		}
	}

	// Prepare revenue data for the last 6 months
	var revenueData []map[string]interface{}
	monthSet := make(map[string]bool) // To keep track of added months

	for i := 0; i < 6; i++ {
		// Calculate revenue for each month
		monthRevenue := 0.0
		monthName := MonthName(i)

		// If the month is already added, skip it to prevent duplicates
		if monthSet[monthName] {
			continue
		}

		// Iterate through invoices to filter by month
		for _, invoice := range allInvoices {
			// Filter by date range of the last 6 months
			if IsInMonth(invoice.Date, i) {
				monthRevenue += invoice.Amount
			}
		}

		// Add month revenue data to the result list
		revenueData = append(revenueData, map[string]interface{}{
			"month":   monthName,
			"revenue": monthRevenue,
		})

		// Mark this month as added
		monthSet[monthName] = true
	}

	// Serialize revenueData into JSON
	revenueDataJSON, err := json.Marshal(revenueData)
	if err != nil {
		logger.Log.Error("Error serializing revenue data", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error serializing revenue data")
	}

	// Log the dashboard data retrieval
	logger.Log.Info("Dashboard data retrieved successfully", zap.Int64("invoiceCount", invoiceCount), zap.Int64("customerCount", customerCount))

	// Check if the request expects JSON (e.g., an AJAX request)
	if c.Query("format") == "json" {
		// Return revenue data as JSON
		return c.JSON(fiber.Map{
			"revenueData": revenueData,
		})
	}

	// Render the dashboard view with the data
	data["InvoiceCount"] = invoiceCount
	data["CustomerCount"] = customerCount
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	data["TotalPaid"] = ac.FormatMoney(int(totalPaid))
	data["TotalPending"] = ac.FormatMoney(int(totalPending))
	data["RevenueData"] = string(revenueDataJSON) // Pass the JSON string to the template
	data["LatestInvoices"] = latestInvoices

	// Render the page
	return c.Render("dashboard", data, "layouts/main")
}

func IsInMonth(invoiceDate string, monthsAgo int) bool {
	date, err := time.Parse("2006-01-02", invoiceDate)
	if err != nil {
		return false
	}

	targetDate := time.Now().AddDate(0, -monthsAgo, 0)
	return date.Year() == targetDate.Year() && date.Month() == targetDate.Month()
}

func MonthName(monthsAgo int) string {
	targetDate := time.Now().AddDate(0, -monthsAgo, 0)
	return targetDate.Format("Jan")
}
