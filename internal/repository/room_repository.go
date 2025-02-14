package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"social_media/internal/domain"
)

type roomRepository struct {
	pool *pgxpool.Pool
}

func NewRoomRepository(pool *pgxpool.Pool) domain.RoomRepository {
	return &roomRepository{pool: pool}
}

func (r *roomRepository) Create(room *domain.Room) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO rooms (id, name, username, type, owner_id, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.pool.Exec(ctx, query,
		room.ID, room.Name, room.Username, room.Type, room.OwnerID, room.CreatedAt, room.UpdatedAt)
	return err
}

func (r *roomRepository) Update(room *domain.Room) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE rooms SET name = $1, username = $2, updated_at = $3 WHERE id = $4`
	_, err := r.pool.Exec(ctx, query, room.Name, room.Username, room.UpdatedAt, room.ID)
	return err
}

func (r *roomRepository) Delete(roomID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, roomID)
	return err
}

func (r *roomRepository) FindByID(roomID string) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, username, type, owner_id, created_at, updated_at FROM rooms WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, roomID)
	var room domain.Room
	err := row.Scan(&room.ID, &room.Name, &room.Username, &room.Type, &room.OwnerID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) FindByUsername(username string) (*domain.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, username, type, owner_id, created_at, updated_at FROM rooms WHERE username = $1`
	row := r.pool.QueryRow(ctx, query, username)
	var room domain.Room
	err := row.Scan(&room.ID, &room.Name, &room.Username, &room.Type, &room.OwnerID, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}
