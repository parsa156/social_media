package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"social_media/internal/domain"
)

type roomMessageRepository struct {
	pool *pgxpool.Pool
}

func NewRoomMessageRepository(pool *pgxpool.Pool) domain.RoomMessageRepository {
	return &roomMessageRepository{pool: pool}
}

func (r *roomMessageRepository) Create(message *domain.RoomMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO room_messages (id, room_id, sender_id, content, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query, message.ID, message.RoomID, message.SenderID, message.Content, message.CreatedAt, message.UpdatedAt)
	return err
}

func (r *roomMessageRepository) Update(message *domain.RoomMessage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE room_messages SET content = $1, updated_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, message.Content, message.UpdatedAt, message.ID)
	return err
}

func (r *roomMessageRepository) Delete(messageID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM room_messages WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, messageID)
	return err
}

func (r *roomMessageRepository) FindByRoom(roomID string) ([]*domain.RoomMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, room_id, sender_id, content, created_at, updated_at FROM room_messages
	          WHERE room_id = $1 ORDER BY created_at ASC`
	rows, err := r.pool.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.RoomMessage
	for rows.Next() {
		var message domain.RoomMessage
		err := rows.Scan(&message.ID, &message.RoomID, &message.SenderID, &message.Content, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (r *roomMessageRepository) FindByID(messageID string) (*domain.RoomMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, room_id, sender_id, content, created_at, updated_at FROM room_messages WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, messageID)
	var message domain.RoomMessage
	err := row.Scan(&message.ID, &message.RoomID, &message.SenderID, &message.Content, &message.CreatedAt, &message.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}
