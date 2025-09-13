package handler

import (
	"net/http"
	"strconv"

	"ecommerce/internal/model"
	"ecommerce/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type TagHandler struct {
	TagService service.ITagService
	Logger     zerolog.Logger
}

func NewTagHandler(service service.ITagService, logger zerolog.Logger) *TagHandler {
	return &TagHandler{
		TagService: service,
		Logger:     logger,
	}
}

// GetAllTags godoc
// @Summary Get all tags
// @Description Retrieve all product tags
// @Tags tags
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/tags [get]
func (h *TagHandler) GetAll(c *gin.Context) {
	tags, err := h.TagService.GetAll(c.Request.Context())
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get tags")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to get tags", err))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Tags retrieved successfully",
		Data:    tags,
	})
}

// GetTagByID godoc
// @Summary Get a tag by ID
// @Description Retrieve a single product tag by its ID
// @Tags tags
// @Produce json
// @Param id path int true "Tag ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/tags/{id} [get]
func (h *TagHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.Logger.Error().Err(err).Msg("invalid tag ID")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid tag ID", err))
		return
	}

	tag, err := h.TagService.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to get tag")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to get tag", err))
		return
	}
	if tag == nil {
		c.JSON(http.StatusNotFound, model.NewErrorResponse(http.StatusNotFound, "Tag not found", nil))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Tag retrieved successfully",
		Data:    tag,
	})
}

// CreateTag godoc
// @Summary Create a new tag
// @Description Create a new product tag (admin only)
// @Tags tags
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param tag body model.CreateTagRequest true "Tag payload"
// @Success 201 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/tags [post]
func (h *TagHandler) Create(c *gin.Context) {
	var req model.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error().Err(err).Msg("invalid request payload")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	tag, err := h.TagService.Create(c.Request.Context(), req)
	if err != nil {
		h.Logger.Error().Err(err).Msg("failed to create tag")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to create tag", err))
		return
	}

	c.JSON(http.StatusCreated, model.SuccessResponse{
		Code:    http.StatusCreated,
		Message: "Tag created successfully",
		Data:    tag,
	})
}

// UpdateTag godoc
// @Summary Update a tag
// @Description Update an existing product tag (admin only)
// @Tags tags
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param id path int true "Tag ID"
// @Param tag body model.UpdateTagRequest true "Update Tag payload"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/tags/{id} [put]
func (h *TagHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.Logger.Error().Err(err).Msg("invalid tag ID")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid tag ID", err))
		return
	}

	var req model.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error().Err(err).Msg("invalid request payload")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid request payload", err))
		return
	}

	tag, err := h.TagService.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		if err.Error() == "tag not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(http.StatusNotFound, "Tag not found", err))
			return
		}
		if err.Error() == "tag name already exists" {
			c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Tag name already exists", err))
			return
		}
		h.Logger.Error().Err(err).Msg("failed to update tag")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to update tag", err))
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{
		Code:    http.StatusOK,
		Message: "Tag updated successfully",
		Data:    tag,
	})
}

// DeleteTag godoc
// @Summary Delete a tag
// @Description Delete a product tag (admin only)
// @Tags tags
// @Security BearerAuth
// @Produce json
// @Param id path int true "Tag ID"
// @Success 204 "No Content"
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/tags/{id} [delete]
func (h *TagHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		h.Logger.Error().Err(err).Msg("invalid tag ID")
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(http.StatusBadRequest, "Invalid tag ID", err))
		return
	}

	err = h.TagService.Delete(c.Request.Context(), uint(id))
	if err != nil {
		if err.Error() == "tag not found" {
			c.JSON(http.StatusNotFound, model.NewErrorResponse(http.StatusNotFound, "Tag not found", err))
			return
		}
		h.Logger.Error().Err(err).Msg("failed to delete tag")
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(http.StatusInternalServerError, "Failed to delete tag", err))
		return
	}

	c.Status(http.StatusNoContent)
}
