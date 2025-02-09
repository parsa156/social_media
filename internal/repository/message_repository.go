package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"social_media/internal/domain"
)

type messageRepository struct {
	pool *pgxpool.Pool
}

func NewMessageRepository(pool *pgxpool.Pool) domain.MessageRepository {
	return &messageRepository{pool: pool}
}

func (r *messageRepository) Create(message *domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO messages (id, conversation_id, sender_id, content, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.pool.Exec(ctx, query,
		message.ID, message.ConversationID, message.SenderID, message.Content, message.CreatedAt, message.UpdatedAt)
	return err
}

func (r *messageRepository) Update(message *domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE messages SET content = $1, updated_at = $2 WHERE id = $3`
	_, err := r.pool.Exec(ctx, query, message.Content, message.UpdatedAt, message.ID)
	return err
}

func (r *messageRepository) Delete(message *domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, message.ID)
	return err
}

func (r *messageRepository) FindByConversation(convoID string) ([]*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, conversation_id, sender_id, content, created_at, updated_at FROM messages
			  WHERE conversation_id = $1 ORDER BY created_at ASC`
	rows, err := r.pool.Query(ctx, query, convoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*domain.Message
	for rows.Next() {
		var message domain.Message
		err := rows.Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.Content, &message.CreatedAt, &message.UpdatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}
	return messages, nil
}

func (r *messageRepository) FindByID(id string) (*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, conversation_id, sender_id, content, created_at, updated_at FROM messages WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	var message domain.Message
	err := row.Scan(&message.ID, &message.ConversationID, &message.SenderID, &message.Content, &message.CreatedAt, &message.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &message, nil
}
