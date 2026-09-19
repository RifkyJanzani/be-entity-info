package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rifky/be-entity-info/internal/model"
	"github.com/rifky/be-entity-info/internal/response"
	"github.com/rifky/be-entity-info/internal/service"
)

// EntityHandler handles HTTP requests for entity endpoints.
type EntityHandler struct {
	service *service.EntityService
}

// NewEntityHandler creates a new EntityHandler.
func NewEntityHandler(svc *service.EntityService) *EntityHandler {
	return &EntityHandler{service: svc}
}

// List handles GET /api/v1/entities
func (h *EntityHandler) List(c *gin.Context) {
	var filter model.EntityFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid query parameters", err.Error())
		return
	}

	entities, err := h.service.ListEntities(c.Request.Context(), filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve entities", nil)
		return
	}

	response.Success(c, http.StatusOK, entities, "Entities retrieved successfully")
}

// GetByID handles GET /api/v1/entities/:id
func (h *EntityHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	entity, err := h.service.GetEntity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Entity %s not found", id), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve entity", nil)
		return
	}

	response.Success(c, http.StatusOK, entity, "Entity retrieved successfully")
}

// Create handles POST /api/v1/entities
func (h *EntityHandler) Create(c *gin.Context) {
	var req model.CreateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", formatValidationErrors(err))
		return
	}

	entity, err := h.service.CreateEntity(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create entity", nil)
		return
	}

	response.Success(c, http.StatusCreated, entity, "Entity created successfully")
}

// Update handles PUT /api/v1/entities/:id
func (h *EntityHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req model.UpdateEntityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", formatValidationErrors(err))
		return
	}

	entity, err := h.service.UpdateEntity(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Entity %s not found", id), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update entity", nil)
		return
	}

	response.Success(c, http.StatusOK, entity, "Entity updated successfully")
}

// Delete handles DELETE /api/v1/entities/:id
func (h *EntityHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.DeleteEntity(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Entity %s not found", id), nil)
			return
		}
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete entity", nil)
		return
	}

	response.Success(c, http.StatusOK, nil, fmt.Sprintf("Entity %s deleted successfully", id))
}

// GetMetrics handles GET /api/v1/entities/metrics
func (h *EntityHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.service.GetMetrics(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to retrieve metrics", nil)
		return
	}

	response.Success(c, http.StatusOK, metrics, "Metrics retrieved successfully")
}
