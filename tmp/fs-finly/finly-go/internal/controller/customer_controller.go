package controller

import (
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

type CustomerController struct {
	customerCollection *mongo.Collection
}

func NewCustomerController(customerCollection *mongo.Collection) *CustomerController {
	return &CustomerController{
		customerCollection: customerCollection,
	}
}

func (cc *CustomerController) Index(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Customers",
		"Type":  "data",
	}

	var customers []model.Customer
	ownerID, err := cc.getOwnerID(c)
	if err != nil {
		data["Customers"] = customers

		logger.Log.Error("Error retrieving owner ID", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	filter := bson.M{
		"owner": ownerID,
	}

	// Check if there is a search query parameter
	search := c.Query("search", "") // Default to empty string if no search is provided
	if search != "" {
		// Build the regex for the search query
		filter["$or"] = []bson.M{
			{"name": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
			{"phone": bson.M{"$regex": search, "$options": "i"}},
			{"address": bson.M{"$regex": search, "$options": "i"}},
		}
	}

	cursor, err := cc.customerCollection.Find(c.Context(), filter)
	if err != nil {
		data["Customers"] = customers

		logger.Log.Error("Error retrieving customers", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}
	defer cursor.Close(c.Context())

	if err := cursor.All(c.Context(), &customers); err != nil {
		data["Customers"] = customers

		logger.Log.Error("Error retrieving customers from cursor", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	logger.Log.Info("Retrieved customers", zap.String("ownerID", ownerID), zap.Int("num_customers", len(customers)))

	data["Customers"] = customers
	data["Search"] = search

	if info, err := sessions.GetInfo(c); err == nil && info != nil {
		data["Info"] = info
	}

	return c.Render("customers/index", data, "layouts/main")
}

func (cc *CustomerController) GetNew(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "New Customer",
		"Type":  "form",
	}
	if info, err := sessions.GetInfo(c); err == nil && info != nil {
		data["Info"] = info
	}

	return c.Render("customers/index", data, "layouts/main")
}

func (cc *CustomerController) PostNew(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "New Customer",
		"Type":  "form",
	}

	var form forms.CustomerNewForm
	if err := c.BodyParser(&form); err != nil {
		data["Error"] = "Error parsing form data"
		data["Customer"] = form

		logger.Log.Error("Error parsing form data", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	if err := form.Validate(); err != nil {
		data["Error"] = err.Error()
		data["Customer"] = form

		logger.Log.Error("Form validation failed", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	profileID, err := sessions.GetProfile(c)
	if err != nil || profileID == "" {
		data["Error"] = "User not authenticated"
		data["Customer"] = form

		logger.Log.Error("User not authenticated", zap.String("profileID", profileID))
		return c.Render("customers/index", data, "layouts/main")
	}

	customer := model.Customer{
		Name:    form.Name,
		Email:   form.Email,
		Phone:   form.Phone,
		Address: form.Address,
		Owner:   profileID,
	}

	if _, err := cc.customerCollection.InsertOne(c.Context(), customer); err != nil {
		data["Error"] = "Error creating customer"
		data["Customer"] = form

		logger.Log.Error("Error creating customer", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	if err := sessions.SetInfo(c, "Customer created successfully", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/dashboard/customers")
}

func (cc *CustomerController) GetEdit(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Edit Customer",
		"Type":  "form",
		"Form":  "edit",
	}

	customerID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		logger.Log.Error("Error converting customer ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID")
	}

	var customer model.Customer
	err = cc.customerCollection.FindOne(c.Context(), bson.M{"_id": objectID}).Decode(&customer)
	if err != nil {
		logger.Log.Error("Error finding customer", zap.Error(err))
		return c.Status(fiber.StatusNotFound).SendString("Customer not found")
	}

	if info, err := sessions.GetInfo(c); err == nil && info != nil {
		data["Info"] = info
	}

	data["Customer"] = customer
	return c.Render("customers/index", data, "layouts/main")
}

func (cc *CustomerController) PostEdit(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Edit Customer",
		"Type":  "form",
		"Form":  "edit",
	}

	var form forms.CustomerEditForm
	if err := c.BodyParser(&form); err != nil {
		data["Error"] = "Error parsing form data"
		data["Customer"] = form

		logger.Log.Error("Error parsing form data", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	if err := form.Validate(); err != nil {
		data["Error"] = err.Error()
		data["Customer"] = form

		logger.Log.Error("Form validation failed", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	updateData := bson.M{
		"name":    form.Name,
		"email":   form.Email,
		"phone":   form.Phone,
		"address": form.Address,
	}

	customerID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		logger.Log.Error("Error converting customer ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID")
	}

	_, err = cc.customerCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": objectID},
		bson.M{"$set": updateData},
	)
	if err != nil {
		data["Error"] = "Error updating customer"
		data["Customer"] = form

		logger.Log.Error("Error updating customer", zap.Error(err))
		return c.Render("customers/index", data, "layouts/main")
	}

	if err := sessions.SetInfo(c, "Customer updated successfully", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/dashboard/customers")
}

func (cc *CustomerController) PostDelete(c *fiber.Ctx) error {
	customerID := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(customerID)
	if err != nil {
		logger.Log.Error("Error converting customer ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID")
	}

	_, err = cc.customerCollection.DeleteOne(c.Context(), bson.M{"_id": objectID})
	if err != nil {
		logger.Log.Error("Error deleting customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting customer")
	}

	if err := sessions.SetInfo(c, "Customer deleted successfully", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/dashboard/customers")
}

func (cc *CustomerController) getOwnerID(c *fiber.Ctx) (string, error) {
	profileID, err := sessions.GetProfile(c)
	if err != nil || profileID == "" {
		return "", err
	}

	return profileID, nil
}
