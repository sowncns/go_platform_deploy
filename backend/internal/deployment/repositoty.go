package deployment

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DeploymentRepository interface {
	Create(ctx context.Context, deployment *Deployment) error
	UpdateStatus(
		ctx context.Context,
		id uint,
		status DeploymentStatus,
	) error
}


type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}


func (r *Repository) Create(ctx context.Context, deployment *Deployment) error {
	query := `
		INSERT INTO deployments (project_id, commit_sha, commit_msg, image_tag, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRow(
		ctx,
		query,
		deployment.ProjectID,
		deployment.CommitSHA,
		deployment.CommitMsg,
		deployment.ImageTag,
		deployment.Status,
	).Scan(&deployment.ID, &deployment.CreatedAt, &deployment.UpdatedAt)
}

func (r *Repository) UpdateStatus(ctx context.Context, id uint, status DeploymentStatus) error {
	query := `
		UPDATE deployments
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := r.db.Exec(ctx, query, status, id)
	return err
}