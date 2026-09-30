package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sowncns/k3s-deploy-platform/internal/auth/dto"
)

var ErrNotFound = errors.New("user not found")

type AuthRepository interface {
	Create(ctx context.Context, req *dto.CreateAuthUser) (*GitHubUser, error)
	getUserGithub(ctx context.Context, userID int64) (*GitHubUser, error)
}

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) getUserGithub(ctx context.Context, userID int64) (*GitHubUser, error) {
	u := &GitHubUser{}
	query := `
		SELECT id
		FROM githubuser
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&u.ID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *Repository) Create(ctx context.Context, req *dto.CreateAuthUser) (*GitHubUser, error) {
	_, err := r.getUserGithub(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	query := `
		INSERT INTO githubuser (
			id,
			login,
			email,
			avatar,
			name
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var user GitHubUser
	err = r.db.QueryRow(
		ctx,
		query,
		req.ID,
		req.Login,
		req.Email,
		req.AvatarURL,
		req.Name,
	).Scan(
		&req.ID,
		&req.Login,
		&req.Email,
		&req.AvatarURL,
		&req.Name,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
