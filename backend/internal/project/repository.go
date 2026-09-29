package project

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("project not found")

type ProjectRepository interface {
	CreateProject(ctx context.Context, proj *Project) error
	UpdateProject(ctx context.Context, proj *Project) error
	DeleteProject(ctx context.Context, proj uint, userID uint) error
	FindProjectByID(ctx context.Context, proj uint) (*Project, error)
	ListProjectsByUserID(ctx context.Context, userID uint) ([]*Project, error)
	PatchEnvVars(ctx context.Context, projectID , userID uint , newVars map[string]string) error
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func marshalEnv(env map[string]string) ([]byte, error) {
	if env == nil {
		env = map[string]string{}
	}
	return json.Marshal(env)
}

func (r *Repository) CreateProject(ctx context.Context, proj *Project) error {
	env, err := marshalEnv(proj.Env)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO projects (
			user_id,
			name,
			git_url,
			branch,
			env
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	return r.db.QueryRow(
		ctx,
		query,
		proj.UserID,
		proj.Name,
		proj.GitURL,
		proj.Branch,
		env,
	).Scan(&proj.ID)
}

func (r *Repository) UpdateProject(ctx context.Context, pro *Project) error {
	env, err := marshalEnv(pro.Env)
	if err != nil {
		return err
	}

	query := `
		UPDATE projects
		SET
			name = $1,
			git_url = $2,
			branch = $3,
			env = $4,
			updated_at = NOW()
		WHERE id = $5 AND user_id = $6
	`

	result, err := r.db.Exec(ctx, query, pro.Name, pro.GitURL, pro.Branch, env, pro.ID, pro.UserID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repository) FindProjectByID(ctx context.Context, id uint) (*Project, error) {
	p := &Project{}
	var env []byte
	query := `
		SELECT id, user_id, name, git_url, branch, env, created_at, updated_at
		FROM projects
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.UserID, &p.Name, &p.GitURL, &p.Branch, &env, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if err := json.Unmarshal(env, &p.Env); err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Repository) ListProjectsByUserID(ctx context.Context, userID uint) ([]*Project, error) {
	query := `
		SELECT id, user_id, name, git_url, branch, env, created_at, updated_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*Project
	for rows.Next() {
		p := &Project{}
		var env []byte
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.GitURL, &p.Branch, &env, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(env, &p.Env); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (r *Repository) DeleteProject(ctx context.Context, projectID uint, userID uint) error {
	query := `
	DELETE FROM projects 
		WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(ctx, query, projectID, userID)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil

}

func (r * Repository) PatchEnvVars(ctx context.Context ,projectID , userID uint, newVars map[string]string) error{
	query := `
        UPDATE projects
        SET 
            env_vars = env_vars || $1::jsonb,
            updated_at = NOW()
        WHERE id = $2 AND user_id = $3
    `
    result, err := r.db.Exec(ctx, query, newVars, projectID, userID)
    if err != nil {
        return err
    }
    if result.RowsAffected() == 0 {
        return pgx.ErrNoRows
    }
    return nil
}