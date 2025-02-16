package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"social_media/internal/domain"
)

type roomMembershipRepository struct {
	pool *pgxpool.Pool
}

func NewRoomMembershipRepository(pool *pgxpool.Pool) domain.RoomMembershipRepository {
	return &roomMembershipRepository{pool: pool}
}

func (r *roomMembershipRepository) AddMember(membership *domain.RoomMembership) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO room_memberships (room_id, user_id, role, created_at)
	          VALUES ($1, $2, $3, $4)`
	_, err := r.pool.Exec(ctx, query, membership.RoomID, membership.UserID, membership.Role, membership.CreatedAt)
	return err
}

func (r *roomMembershipRepository) UpdateMemberRole(roomID, userID string, role domain.RoomMembershipRole) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `UPDATE room_memberships SET role = $1 WHERE room_id = $2 AND user_id = $3`
	cmdTag, err := r.pool.Exec(ctx, query, role, roomID, userID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("membership not found")
	}
	return nil
}

func (r *roomMembershipRepository) RemoveMember(roomID, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `DELETE FROM room_memberships WHERE room_id = $1 AND user_id = $2`
	_, err := r.pool.Exec(ctx, query, roomID, userID)
	return err
}

func (r *roomMembershipRepository) GetMembers(roomID string) ([]*domain.RoomMembership, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT room_id, user_id, role, created_at FROM room_memberships WHERE room_id = $1`
	rows, err := r.pool.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var memberships []*domain.RoomMembership
	for rows.Next() {
		var m domain.RoomMembership
		err := rows.Scan(&m.RoomID, &m.UserID, &m.Role, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, &m)
	}
	return memberships, nil
}

func (r *roomMembershipRepository) IsUserBanned(roomID, userID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT role FROM room_memberships WHERE room_id = $1 AND user_id = $2`
	row := r.pool.QueryRow(ctx, query, roomID, userID)
	var role string
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return role == string(domain.RoleBanned), nil
}

func (r *roomMembershipRepository) GetMemberRole(roomID, userID string) (domain.RoomMembershipRole, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT role FROM room_memberships WHERE room_id = $1 AND user_id = $2`
	row := r.pool.QueryRow(ctx, query, roomID, userID)
	var role string
	err := row.Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return domain.RoomMembershipRole(role), nil
}
