package repository

import (
	"context"
	"errors"
	"time"

	"traffic-control/account/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type IUserRepository interface {
	FindAll(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type userRepository struct {
	pool *pgxpool.Pool
	log  zerolog.Logger
}

func NewUserRepository(pool *pgxpool.Pool, log zerolog.Logger) IUserRepository {
	return &userRepository{
		pool: pool,
		log:  log.With().Str("component", "repository.user").Logger(),
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	now := time.Now().UTC().Format(time.RFC3339)
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.pool.Exec(ctx, query, user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		r.log.Error().Err(err).Str("email", user.Email).Msg("Failed to create user")
		return err
	}

	r.log.Info().Str("id", user.ID).Str("email", user.Email).Msg("User created successfully")
	return nil
}

func (r *userRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		r.log.Error().Err(err).Msg("Failed to query all users")
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
			r.log.Error().Err(err).Msg("Failed to scan user row")
			continue
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.log.Error().Err(err).Msg("Error iterating user rows")
		return nil, err
	}

	r.log.Debug().Int("count", len(users)).Msg("Users retrieved")
	return users, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)
	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn().Str("id", id).Msg("User not found")
			return nil, nil
		}
		r.log.Error().Err(err).Str("id", id).Msg("Failed to find user by ID")
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, updated_at = $4
		WHERE id = $5
	`

	cmd, err := r.pool.Exec(ctx, query, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		r.log.Error().Err(err).Str("id", user.ID).Msg("Failed to update user")
		return err
	}

	if cmd.RowsAffected() == 0 {
		r.log.Warn().Str("id", user.ID).Msg("No user updated (not found)")
		return nil
	}

	r.log.Info().Str("id", user.ID).Msg("User updated successfully")
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`

	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		r.log.Error().Err(err).Str("id", id).Msg("Failed to delete user")
		return err
	}

	if cmd.RowsAffected() == 0 {
		r.log.Warn().Str("id", id).Msg("No user deleted (not found)")
		return nil
	}

	r.log.Info().Str("id", id).Msg("User deleted successfully")
	return nil
}
