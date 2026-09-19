package model

import "time"

// Entity represents a geospatial entity stored in the database.
type Entity struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Kind        string      `json:"kind"`
	Status      string      `json:"status"`
	Latitude    float64     `json:"latitude"`
	Longitude   float64     `json:"longitude"`
	Description *string     `json:"description"`
	Attributes  []Attribute `json:"attributes"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// Attribute represents a key-value pair attached to an entity.
type Attribute struct {
	Label string `json:"label" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// MetricsResponse holds aggregated dashboard metrics.
type MetricsResponse struct {
	TotalCount    int `json:"totalCount"`
	ActiveCount   int `json:"activeCount"`
	OfflineCount  int `json:"offlineCount"`
	FacilityCount int `json:"facilityCount"`
}
