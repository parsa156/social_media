package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"social_media/internal/domain"
)

type conversationRepository struct {
	pool *pgxpool.Pool
}

func NewConversationRepository(pool *pgxpool.Pool) domain.ConversationRepository {
	return &conversationRepository{pool: pool}
}

func (r *conversationRepository) Create(convo *domain.Conversation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO conversations (id, participant1, participant2, created_at)
			  VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, convo.ID, convo.Participant1, convo.Participant2, convo.CreatedAt)
	return err
}

func (r *conversationRepository) FindByParticipants(p1, p2 string) (*domain.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, participant1, participant2, created_at FROM conversations
			  WHERE participant1 = $1 AND participant2 = $2`
	row := r.pool.QueryRow(ctx, query, p1, p2)
	var convo domain.Conversation
	err := row.Scan(&convo.ID, &convo.Participant1, &convo.Participant2, &convo.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &convo, nil
}

func (r *conversationRepository) FindByUser(userID string) ([]*domain.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, participant1, participant2, created_at FROM conversations
			  WHERE participant1 = $1 OR participant2 = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var convos []*domain.Conversation
	for rows.Next() {
		var convo domain.Conversation
		err := rows.Scan(&convo.ID, &convo.Participant1, &convo.Participant2, &convo.CreatedAt)
		if err != nil {
			return nil, err
		}
		convos = append(convos, &convo)
	}
	return convos, nil
}

func (r *conversationRepository) FindByID(id string) (*domain.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, participant1, participant2, created_at FROM conversations WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)
	var convo domain.Conversation
	err := row.Scan(&convo.ID, &convo.Participant1, &convo.Participant2, &convo.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &convo, nil
}

