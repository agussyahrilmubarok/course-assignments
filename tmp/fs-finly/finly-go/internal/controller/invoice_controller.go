package controller

import (
	"context"
	"finly-go/internal/forms"
	"finly-go/internal/model"
	"finly-go/internal/sessions"
	"finly-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type InvoiceController struct {
	invoiceCollection  *mongo.Collection
	customerCollection *mongo.Collection
}

func NewInvoiceController(invoiceCollection, customerCollection *mongo.Collection) *InvoiceController {
	return &InvoiceController{
		invoiceCollection:  invoiceCollection,
		customerCollection: customerCollection,
	}
}

func (ic *InvoiceController) Index(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Invoices",
		"Type":  "data",
	}

	var invoices []model.Invoice
	ownerID, err := ic.getOwnerID(c)
	if err != nil {
		logger.Log.Error("Error retrieving owner ID", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	filter := bson.M{
		"owner": ownerID,
	}

	search := c.Query("search", "")
	if search != "" {
		filter["$or"] = []bson.M{
			{"customer": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	cursor, err := ic.invoiceCollection.Find(c.Context(), filter)
	if err != nil {
		logger.Log.Error("Error retrieving invoices", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}
	defer cursor.Close(c.Context())

	if err := cursor.All(c.Context(), &invoices); err != nil {
		logger.Log.Error("Error retrieving invoices from cursor", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	customerIDs := make([]string, len(invoices))
	for i, invoice := range invoices {
		customerIDs[i] = invoice.Customer
	}

	var customers []model.Customer
	if len(customerIDs) > 0 {
		customerFilter := bson.M{
			"_id": bson.M{"$in": make([]interface{}, len(customerIDs))},
		}
		for i, customerID := range customerIDs {
			objectID, err := primitive.ObjectIDFromHex(customerID)
			if err != nil {
				logger.Log.Error("Error converting string to ObjectId", zap.String("customerID", customerID), zap.Error(err))
				return c.Render("invoices/index", data, "layouts/main")
			}
			customerFilter["_id"].(bson.M)["$in"].([]interface{})[i] = objectID
		}

		logger.Log.Info("Querying customers", zap.Any("customerFilter", customerFilter))

		cursor, err := ic.customerCollection.Find(c.Context(), customerFilter)
		if err != nil {
			logger.Log.Error("Error retrieving customers", zap.Error(err))
			return c.Render("invoices/index", data, "layouts/main")
		}
		defer cursor.Close(c.Context())

		if err := cursor.All(c.Context(), &customers); err != nil {
			logger.Log.Error("Error retrieving customers from cursor", zap.Error(err))
			return c.Render("invoices/index", data, "layouts/main")
		}
	}

	customerMap := make(map[string]string)
	for _, customer := range customers {
		customerMap[customer.ID] = customer.Name
	}

	for i := range invoices {
		if name, ok := customerMap[invoices[i].Customer]; ok {
			invoices[i].CustomerName = name
		}
	}

	data["Invoices"] = invoices
	data["Search"] = search

	if info, err := sessions.GetInfo(c); err == nil && info != nil {
		data["Info"] = info
	}

	return c.Render("invoices/index", data, "layouts/main")
}

func (ic *InvoiceController) GetNew(c *fiber.Ctx) error {
	customers := ic.getCustomersByOwnerID(c, ic.customerCollection)
	data := fiber.Map{
		"Title":     "New Invoice",
		"Type":      "form",
		"Customers": customers,
	}

	if info, err := sessions.GetInfo(c); err == nil && info != nil {
		data["Info"] = info
	}

	return c.Render("invoices/index", data, "layouts/main")
}

func (ic *InvoiceController) PostNew(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "New Invoice",
		"Type":  "form",
	}

	var form forms.InvoiceNewForm
	if err := c.BodyParser(&form); err != nil {
		data["Customers"] = ic.getCustomersByOwnerID(c, ic.customerCollection)
		data["Error"] = "Error parsing form data"
		data["Invoice"] = form

		logger.Log.Error("Error parsing form data", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	if err := form.Validate(); err != nil {
		data["Customers"] = ic.getCustomersByOwnerID(c, ic.customerCollection)
		data["Error"] = err.Error()
		data["Invoice"] = form

		logger.Log.Error("Form validation failed", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	profileID, err := sessions.GetProfile(c)
	if err != nil || profileID == "" {
		data["Customers"] = ic.getCustomersByOwnerID(c, ic.customerCollection)
		data["Error"] = "User not authenticated"
		data["Invoice"] = form

		logger.Log.Error("User not authenticated", zap.String("profileID", profileID))
		return c.Render("invoices/index", data, "layouts/main")
	}

	invoice := model.Invoice{
		Amount:   form.Amount,
		Date:     form.Date,
		Status:   form.Status,
		Owner:    profileID,
		Customer: form.Customer,
	}

	if _, err := ic.invoiceCollection.InsertOne(c.Context(), invoice); err != nil {
		data["Customers"] = ic.getCustomersByOwnerID(c, ic.customerCollection)
		data["Error"] = "Error creating invoice"
		data["Invoice"] = form

		logger.Log.Error("Error creating invoice", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	if err := sessions.SetInfo(c, "Invoice created successfully", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/dashboard/invoices")
}

func (ic *InvoiceController) GetEdit(c *fiber.Ctx) error {
	customers := ic.getCustomersByOwnerID(c, ic.customerCollection)
	data := fiber.Map{
		"Title":     "Edit Invoice",
		"Type":      "form",
		"Form":      "edit",
		"Customers": customers,
	}

	invoiceID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		logger.Log.Error("Error converting invoice ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid invoice ID")
	}

	var invoice model.Invoice
	err = ic.invoiceCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&invoice)
	if err != nil {
		logger.Log.Error("Error finding invoice", zap.Error(err))
		return c.Status(fiber.StatusNotFound).SendString("Invoice not found")
	}

	if info, err := sessions.GetInfo(c); err == nil && info != nil {
		data["Info"] = info
	}

	data["Invoice"] = invoice
	return c.Render("invoices/index", data, "layouts/main")
}

func (ic *InvoiceController) PostEdit(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Edit Invoice",
		"Type":  "form",
		"Form":  "edit",
	}

	var form forms.InvoiceEditForm
	if err := c.BodyParser(&form); err != nil {
		data["Error"] = "Error parsing form data"
		data["Invoice"] = form

		logger.Log.Error("Error parsing form data", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	if err := form.Validate(); err != nil {
		data["Error"] = err.Error()
		data["Invoice"] = form

		logger.Log.Error("Form validation failed", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	updateData := bson.M{
		"amount":   form.Amount,
		"date":     form.Date,
		"status":   form.Status,
		"customer": form.Customer,
	}

	invoiceID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(invoiceID)
	if err != nil {
		logger.Log.Error("Error converting invoice ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid invoice ID")
	}

	_, err = ic.invoiceCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": objectID},
		bson.M{"$set": updateData},
	)
	if err != nil {
		data["Error"] = "Error updating invoice"
		data["Invoice"] = form

		logger.Log.Error("Error updating invoice", zap.Error(err))
		return c.Render("invoices/index", data, "layouts/main")
	}

	if err := sessions.SetInfo(c, "Invoice updated successfully", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/dashboard/invoices")
}

func (cc *InvoiceController) PostDelete(c *fiber.Ctx) error {
	customerID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		logger.Log.Error("Error converting customer ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID")
	}

	_, err = cc.invoiceCollection.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		logger.Log.Error("Error deleting invoice", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting invoice")
	}

	if err := sessions.SetInfo(c, "Invoice deleted successfully", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/dashboard/invoices")
}

func (ic *InvoiceController) getOwnerID(c *fiber.Ctx) (string, error) {
	profileID, err := sessions.GetProfile(c)
	if err != nil || profileID == "" {
		logger.Log.Error("Error retrieving owner ID", zap.String("profileID", profileID))
		return "", err
	}
	return profileID, nil
}

func (ic *InvoiceController) getCustomersByOwnerID(c *fiber.Ctx, customerCollection *mongo.Collection) []model.Customer {
	var customers []model.Customer
	ownerID, err := ic.getOwnerID(c)
	if err != nil {
		logger.Log.Error("Error retrieving owner ID", zap.Error(err))
		return customers
	}

	cursor, err := customerCollection.Find(context.Background(), bson.M{"owner": ownerID})
	if err != nil {
		logger.Log.Error("Error retrieving customers", zap.Error(err))
		return customers
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &customers); err != nil {
		logger.Log.Error("Error retrieving customers from cursor", zap.Error(err))
		return customers
	}

	return customers
}

func (ic *InvoiceController) getCustomerByID(c *fiber.Ctx, customerID string) string {
	var customer model.Customer
	err := ic.customerCollection.FindOne(c.Context(), bson.M{"_id": customerID}).Decode(&customer)
	if err != nil {
		return customer.Name
	}
	return customer.Name
}
