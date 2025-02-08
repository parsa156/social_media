package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"social_media/internal/domain"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) domain.UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users (id, name, phone, username, password, created_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Name, user.Phone, user.Username, user.Password, user.CreatedAt)
	return err
}

func (r *userRepository) FindByPhone(phone string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, phone, username, password, created_at FROM users WHERE phone = $1`
	row := r.pool.QueryRow(ctx, query, phone)
	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, phone, username, password, created_at FROM users WHERE username = $1`
	row := r.pool.QueryRow(ctx, query, username)
	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, phone, username, password, created_at FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	var user domain.User
	err := row.Scan(&user.ID, &user.Name, &user.Phone, &user.Username, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE users SET name = $1, phone = $2, username = $3, password = $4 WHERE id = $5`
	_, err := r.pool.Exec(ctx, query, user.Name, user.Phone, user.Username, user.Password, user.ID)
	return err
}

func (r *userRepository) Delete(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM users WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, user.ID)
	return err
}
