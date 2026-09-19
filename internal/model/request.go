package model

// CreateEntityRequest is the input DTO for creating a new entity.
type CreateEntityRequest struct {
	Name        string      `json:"name" binding:"required,min=1,max=255"`
	Kind        string      `json:"kind" binding:"required,entity_kind"`
	Status      string      `json:"status" binding:"required,entity_status"`
	Latitude    *float64    `json:"latitude" binding:"required,min=-90,max=90"`
	Longitude   *float64    `json:"longitude" binding:"required,min=-180,max=180"`
	Description *string     `json:"description" binding:"omitempty,max=1000"`
	Attributes  []Attribute `json:"attributes" binding:"omitempty,dive"`
}

// UpdateEntityRequest is the input DTO for updating an existing entity (full replace).
type UpdateEntityRequest = CreateEntityRequest

// EntityFilter holds query parameters for filtering the entity list.
type EntityFilter struct {
	Kind   string `form:"kind"`
	Status string `form:"status"`
	Search string `form:"search"`
}
