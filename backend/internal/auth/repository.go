package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository interface {
	FindByID(ctx context.Context, id int64) (*GitHubUser, error)
	Upsert(ctx context.Context, user *GitHubUser) error
	FindByAccessToken(ctx context.Context, accessToken string) (*GitHubUser, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) AuthRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) FindByID(
	ctx context.Context,
	id int64,
) (*GitHubUser, error) {
	user := &GitHubUser{}

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			login,
			name,
			email,
			avatar_url,
			access_token,
			created_at,
			updated_at
		FROM github_users
		WHERE id = $1
		`,
		id,
	).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.Email,
		&user.AvatarURL,
		&user.AccessToken,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("find github user by id: %w", err)
	}

	return user, nil
}

// Upsert creates the user if it doesn't exist (by GitHub id), or updates
// its profile fields and access token if it does.
func (r *repository) Upsert(
	ctx context.Context,
	user *GitHubUser,
) error {
	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO github_users (
			id,
			login,
			name,
			email,
			avatar_url,
			access_token
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			login = EXCLUDED.login,
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			avatar_url = EXCLUDED.avatar_url,
			access_token = EXCLUDED.access_token,
			updated_at = NOW()
		RETURNING created_at, updated_at
		`,
		user.ID,
		user.Login,
		user.Name,
		user.Email,
		user.AvatarURL,
		user.AccessToken,
	).Scan(
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("upsert github user: %w", err)
	}

	return nil
}


func (r *repository) FindByAccessToken(ctx context.Context, accessToken string) (*GitHubUser, error) {
    var user GitHubUser
   query := `SELECT id, login, name, email, avatar_url, access_token FROM github_users WHERE access_token = $1`
   err := r.db.QueryRow(ctx, query, accessToken).Scan(
		&user.ID,
		&user.Login,
		&user.Name,
		&user.Email,
		&user.AvatarURL,
		&user.AccessToken,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Không tìm thấy user
		}
		return nil, err // Lỗi database khác
	}

	return &user, nil
}