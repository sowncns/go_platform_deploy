package cluster

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("cluster not found")

type ClusterRepository interface {
	Create(ctx context.Context, cluster *Cluster) error
	UpdateStatus(ctx context.Context, id string, status ClusterStatus) error
	FindByID(ctx context.Context, id string) (*Cluster, error)
	List(ctx context.Context) ([]*Cluster, error)
	Delete(ctx context.Context, id string) error
}


type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, c *Cluster) error {
	query := `
		INSERT INTO clusters (
			id, name, provider, status, endpoint, region, zone, credential_ref, node_count
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	return r.db.QueryRow(
		ctx, query,
		c.ID, c.Name, c.Provider, c.Status, c.Endpoint, c.Region, c.Zone, c.CredentialRef, c.NodeCount,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
}

func (r *Repository) UpdateStatus(ctx context.Context, id string, status ClusterStatus) error {
	query := `UPDATE clusters SET status = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (*Cluster, error) {
	query := `
		SELECT id, name, provider, status, endpoint, region, zone, credential_ref, node_count, created_at, updated_at
		FROM clusters
		WHERE id = $1
	`

	c := &Cluster{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.Name, &c.Provider, &c.Status, &c.Endpoint, &c.Region, &c.Zone,
		&c.CredentialRef, &c.NodeCount, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return c, nil
}

func (r *Repository) List(ctx context.Context) ([]*Cluster, error) {
	query := `
		SELECT id, name, provider, status, endpoint, region, zone, credential_ref, node_count, created_at, updated_at
		FROM clusters
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clusters []*Cluster
	for rows.Next() {
		c := &Cluster{}
		if err := rows.Scan(
			&c.ID, &c.Name, &c.Provider, &c.Status, &c.Endpoint, &c.Region, &c.Zone,
			&c.CredentialRef, &c.NodeCount, &c.CreatedAt, &c.UpdatedAt,
		); err != nil {
			return nil, err
		}
		clusters = append(clusters, c)
	}
	return clusters, rows.Err()
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	result, err := r.db.Exec(ctx, `DELETE FROM clusters WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
