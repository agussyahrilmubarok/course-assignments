package handler

import (
	"net/http"

	"example.com/pkg/model"
	"example.com/productfc/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type productHandler struct {
	productService service.IProductService
	log            *zap.Logger
}

func NewProductHandler(productService service.IProductService, log *zap.Logger) *productHandler {
	return &productHandler{
		productService: productService,
		log:            log,
	}
}

func (h *productHandler) CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid create product request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	response, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil || response == nil {
		h.log.Error("failed to create product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed create product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *productHandler) UpdateProductByID(c *gin.Context) {
	id := c.Param("id")
	var req model.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("invalid update product request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	response, err := h.productService.UpdateByID(c.Request.Context(), id, &req)
	if err != nil || response == nil {
		h.log.Error("failed to update product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *productHandler) DeleteProductByID(c *gin.Context) {
	id := c.Param("id")
	if err := h.productService.DeleteByID(c.Request.Context(), id); err != nil {
		h.log.Error("failed to delete product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (h *productHandler) FindProductByID(c *gin.Context) {
	id := c.Param("id")
	response, err := h.productService.FindByID(c.Request.Context(), id)
	if err != nil || response == nil {
		h.log.Error("failed to find product", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find product", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *productHandler) SearchProduct(c *gin.Context) {

}
