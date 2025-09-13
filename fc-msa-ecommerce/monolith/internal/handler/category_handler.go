package handler

import (
	"net/http"
	"strconv"

	"ecommerce/internal/model"
	"ecommerce/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type CategoryHandler struct {
	CategoryService service.ICategoryService
	Logger          zerolog.Logger
}

func NewCategoryHandler(service service.ICategoryService, logger zerolog.Logger) *CategoryHandler {
	return &CategoryHandler{
		CategoryService: service,
		Logger:          logger,
	}
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve all product categories
// @Tags categories
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.CategoryService.GetAll(c.Request.Context())
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get categories")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to get categories", err))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

// GetCategoryByID godoc
// @Summary Get a category by ID
// @Description Retrieve a single product category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.Logger.Error().Err(err).Msg("invalid category ID")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid category ID", err))
		return
	}

	category, err := h.CategoryService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get category")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to get category", err))
		return
	}
	if category == nil {
		c.JSON(http.StatusNotFound, model.NewErrorResponse(http.StatusNotFound, "Category not found", nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Category retrieved successfully",
		Data:    category,
	})
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param category body model.CreateCategoryRequest true "Category payload"
// @Success 201 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/categories [post]
func (h *CategoryHandler) Create(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error().Err(err).Msg("invalid request payload")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	category, err := h.CategoryService.Create(c.Request.Context(), req)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to create category")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to create category", err))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "Category created successfully",
		Data:    category,
	})
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing product category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body model.UpdateCategoryRequest true "Update Category payload"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.Logger.Error().Err(err).Msg("invalid category ID")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid category ID", err))
		return
	}

	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error().Err(err).Msg("invalid request payload")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	category, err := h.CategoryService.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(http.StatusNotFound, "Category not found", err))
			return
		}
		if err.Error() == "category name already exists" {
			c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Category name already exists", err))
			return
		}
		h.Logger.Error().Err(err).Msg("failed to update category")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to update category", err))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Category updated successfully",
		Data:    category,
	})
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a product category (admin only)
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 204 "No Content"
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.Logger.Error().Err(err).Msg("invalid category ID")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid category ID", err))
		return
	}

	err = h.CategoryService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "category not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(http.StatusNotFound, "Category not found", err))
			return
		}
		h.Logger.Error().Err(err).Msg("failed to delete category")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to delete category", err))
		return
	}

	c.Status(http.StatusNoContent)
}
