package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rifky/be-entity-info/internal/model"
)

// psql is a squirrel statement builder configured for PostgreSQL ($1, $2, ...) placeholders.
var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

// EntityRepository handles all database operations for entities.
type EntityRepository struct {
	pool *pgxpool.Pool
}

// NewEntityRepository creates a new EntityRepository.
func NewEntityRepository(pool *pgxpool.Pool) *EntityRepository {
	return &EntityRepository{pool: pool}
}

// List retrieves all entities with optional filters, including their attributes.
func (r *EntityRepository) List(ctx context.Context, filter model.EntityFilter) ([]model.Entity, error) {
	qb := psql.Select(
		"e.id", "e.name", "e.kind", "e.status",
		"e.latitude", "e.longitude", "e.description",
		"e.created_at", "e.updated_at",
	).From("entities e")

	if filter.Kind != "" {
		qb = qb.Where(sq.Eq{"e.kind": filter.Kind})
	}
	if filter.Status != "" {
		qb = qb.Where(sq.Eq{"e.status": filter.Status})
	}
	if filter.Search != "" {
		qb = qb.Where(sq.ILike{"e.name": fmt.Sprintf("%%%s%%", filter.Search)})
	}

	qb = qb.OrderBy("e.updated_at DESC")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build list query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query entities: %w", err)
	}
	defer rows.Close()

	var entities []model.Entity
	for rows.Next() {
		var e model.Entity
		if err := rows.Scan(
			&e.ID, &e.Name, &e.Kind, &e.Status,
			&e.Latitude, &e.Longitude, &e.Description,
			&e.CreatedAt, &e.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan entity: %w", err)
		}
		entities = append(entities, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate entities: %w", err)
	}

	// Load attributes for all entities in one query
	if len(entities) > 0 {
		if err := r.loadAttributes(ctx, entities); err != nil {
			return nil, err
		}
	}

	return entities, nil
}

// GetByID retrieves a single entity by ID, including its attributes.
func (r *EntityRepository) GetByID(ctx context.Context, id string) (*model.Entity, error) {
	query, args, err := psql.Select(
		"id", "name", "kind", "status",
		"latitude", "longitude", "description",
		"created_at", "updated_at",
	).From("entities").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get query: %w", err)
	}

	var e model.Entity
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&e.ID, &e.Name, &e.Kind, &e.Status,
		&e.Latitude, &e.Longitude, &e.Description,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query entity by id: %w", err)
	}

	// Load attributes
	attrs, err := r.getAttributesByEntityID(ctx, id)
	if err != nil {
		return nil, err
	}
	e.Attributes = attrs

	return &e, nil
}

// Create inserts a new entity with its attributes in a transaction.
func (r *EntityRepository) Create(ctx context.Context, entity *model.Entity) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert entity
	query, args, err := psql.Insert("entities").
		Columns("id", "name", "kind", "status", "latitude", "longitude", "description", "created_at", "updated_at").
		Values(entity.ID, entity.Name, entity.Kind, entity.Status, entity.Latitude, entity.Longitude, entity.Description, entity.CreatedAt, entity.UpdatedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("build insert query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert entity: %w", err)
	}

	// Insert attributes
	if err := r.insertAttributes(ctx, tx, entity.ID, entity.Attributes); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Update modifies an existing entity and replaces its attributes in a transaction.
func (r *EntityRepository) Update(ctx context.Context, entity *model.Entity) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update entity
	query, args, err := psql.Update("entities").
		Set("name", entity.Name).
		Set("kind", entity.Kind).
		Set("status", entity.Status).
		Set("latitude", entity.Latitude).
		Set("longitude", entity.Longitude).
		Set("description", entity.Description).
		Set("updated_at", entity.UpdatedAt).
		Where(sq.Eq{"id": entity.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build update query: %w", err)
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update entity: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	// Delete old attributes, insert new ones
	delQuery, delArgs, _ := psql.Delete("entity_attributes").Where(sq.Eq{"entity_id": entity.ID}).ToSql()
	if _, err := tx.Exec(ctx, delQuery, delArgs...); err != nil {
		return fmt.Errorf("delete old attributes: %w", err)
	}

	if err := r.insertAttributes(ctx, tx, entity.ID, entity.Attributes); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Delete removes an entity by ID. Attributes are cascade-deleted by the database.
func (r *EntityRepository) Delete(ctx context.Context, id string) error {
	query, args, err := psql.Delete("entities").Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return fmt.Errorf("build delete query: %w", err)
	}

	result, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete entity: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// GetMetrics returns aggregated counts for the dashboard metrics cards.
func (r *EntityRepository) GetMetrics(ctx context.Context) (*model.MetricsResponse, error) {
	query := `
		SELECT
			COUNT(*) AS total_count,
			COUNT(*) FILTER (WHERE status = 'Active') AS active_count,
			COUNT(*) FILTER (WHERE status = 'Offline') AS offline_count,
			COUNT(*) FILTER (WHERE kind = 'Facility') AS facility_count
		FROM entities
	`

	var m model.MetricsResponse
	err := r.pool.QueryRow(ctx, query).Scan(
		&m.TotalCount, &m.ActiveCount, &m.OfflineCount, &m.FacilityCount,
	)
	if err != nil {
		return nil, fmt.Errorf("query metrics: %w", err)
	}

	return &m, nil
}

// loadAttributes fetches attributes for a slice of entities in a single query.
func (r *EntityRepository) loadAttributes(ctx context.Context, entities []model.Entity) error {
	ids := make([]string, len(entities))
	for i, e := range entities {
		ids[i] = e.ID
	}

	query, args, err := psql.Select("entity_id", "label", "value").
		From("entity_attributes").
		Where(sq.Eq{"entity_id": ids}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return fmt.Errorf("build attributes query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query attributes: %w", err)
	}
	defer rows.Close()

	attrMap := make(map[string][]model.Attribute)
	for rows.Next() {
		var entityID string
		var attr model.Attribute
		if err := rows.Scan(&entityID, &attr.Label, &attr.Value); err != nil {
			return fmt.Errorf("scan attribute: %w", err)
		}
		attrMap[entityID] = append(attrMap[entityID], attr)
	}

	for i := range entities {
		if attrs, ok := attrMap[entities[i].ID]; ok {
			entities[i].Attributes = attrs
		} else {
			entities[i].Attributes = []model.Attribute{}
		}
	}

	return nil
}

// getAttributesByEntityID fetches attributes for a single entity.
func (r *EntityRepository) getAttributesByEntityID(ctx context.Context, entityID string) ([]model.Attribute, error) {
	query, args, err := psql.Select("label", "value").
		From("entity_attributes").
		Where(sq.Eq{"entity_id": entityID}).
		OrderBy("id ASC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build attributes query: %w", err)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query attributes: %w", err)
	}
	defer rows.Close()

	var attrs []model.Attribute
	for rows.Next() {
		var a model.Attribute
		if err := rows.Scan(&a.Label, &a.Value); err != nil {
			return nil, fmt.Errorf("scan attribute: %w", err)
		}
		attrs = append(attrs, a)
	}

	if attrs == nil {
		attrs = []model.Attribute{}
	}

	return attrs, nil
}

// insertAttributes batch-inserts attributes for an entity within a transaction.
func (r *EntityRepository) insertAttributes(ctx context.Context, tx pgx.Tx, entityID string, attrs []model.Attribute) error {
	if len(attrs) == 0 {
		return nil
	}

	qb := psql.Insert("entity_attributes").Columns("entity_id", "label", "value")
	for _, a := range attrs {
		qb = qb.Values(entityID, a.Label, a.Value)
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return fmt.Errorf("build insert attributes query: %w", err)
	}

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("insert attributes: %w", err)
	}

	return nil
}
