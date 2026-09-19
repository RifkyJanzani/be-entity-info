package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rifky/be-entity-info/internal/model"
	"github.com/rifky/be-entity-info/internal/repository"
)

// EntityService contains business logic for entity management.
type EntityService struct {
	repo *repository.EntityRepository
}

// NewEntityService creates a new EntityService.
func NewEntityService(repo *repository.EntityRepository) *EntityService {
	return &EntityService{repo: repo}
}

// ListEntities returns all entities matching the given filter.
func (s *EntityService) ListEntities(ctx context.Context, filter model.EntityFilter) ([]model.Entity, error) {
	entities, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list entities: %w", err)
	}

	// Ensure non-nil slice for consistent JSON response
	if entities == nil {
		entities = []model.Entity{}
	}

	return entities, nil
}

// GetEntity returns a single entity by ID.
func (s *EntityService) GetEntity(ctx context.Context, id string) (*model.Entity, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get entity: %w", err)
	}
	if entity == nil {
		return nil, ErrNotFound
	}

	return entity, nil
}

// CreateEntity validates the request, generates an ID, and persists a new entity.
func (s *EntityService) CreateEntity(ctx context.Context, req model.CreateEntityRequest) (*model.Entity, error) {
	now := time.Now()

	entity := &model.Entity{
		ID:          generateID(),
		Name:        strings.TrimSpace(req.Name),
		Kind:        req.Kind,
		Status:      req.Status,
		Latitude:    *req.Latitude,
		Longitude:   *req.Longitude,
		Description: trimStringPtr(req.Description),
		Attributes:  sanitizeAttributes(req.Attributes),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf("create entity: %w", err)
	}

	return entity, nil
}

// UpdateEntity validates the request and updates an existing entity.
func (s *EntityService) UpdateEntity(ctx context.Context, id string, req model.UpdateEntityRequest) (*model.Entity, error) {
	// Check existence first
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("check entity existence: %w", err)
	}
	if existing == nil {
		return nil, ErrNotFound
	}

	entity := &model.Entity{
		ID:          id,
		Name:        strings.TrimSpace(req.Name),
		Kind:        req.Kind,
		Status:      req.Status,
		Latitude:    *req.Latitude,
		Longitude:   *req.Longitude,
		Description: trimStringPtr(req.Description),
		Attributes:  sanitizeAttributes(req.Attributes),
		CreatedAt:   existing.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Update(ctx, entity); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update entity: %w", err)
	}

	return entity, nil
}

// DeleteEntity removes an entity by ID.
func (s *EntityService) DeleteEntity(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf("delete entity: %w", err)
	}
	return nil
}

// GetMetrics returns aggregated dashboard metrics.
func (s *EntityService) GetMetrics(ctx context.Context) (*model.MetricsResponse, error) {
	return s.repo.GetMetrics(ctx)
}

// generateID creates a unique entity ID in the format ENT-XXXXXX (6 hex chars).
func generateID() string {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	return fmt.Sprintf("ENT-%X", b)
}

// sanitizeAttributes trims whitespace from attribute labels and values.
func sanitizeAttributes(attrs []model.Attribute) []model.Attribute {
	if attrs == nil {
		return []model.Attribute{}
	}
	result := make([]model.Attribute, len(attrs))
	for i, a := range attrs {
		result[i] = model.Attribute{
			Label: strings.TrimSpace(a.Label),
			Value: strings.TrimSpace(a.Value),
		}
	}
	return result
}

// trimStringPtr trims whitespace from a string pointer.
func trimStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
