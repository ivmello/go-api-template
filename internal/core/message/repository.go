package message

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrMessageNotFound = errors.New("message not found")
	ErrForbidden       = errors.New("access forbidden")
)

// Repository provides access to message storage
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new message repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

// Create inserts a new message into the database
func (r *Repository) Create(ctx context.Context, message *Message) error {
	query := `
		INSERT INTO messages (user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	
	return r.db.QueryRow(ctx, query,
		message.UserID,
		message.Content,
		message.CreatedAt,
		message.UpdatedAt,
	).Scan(&message.ID)
}

// GetAll retrieves all messages
func (r *Repository) GetAll(ctx context.Context) ([]*Message, error) {
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM messages
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*Message
	for rows.Next() {
		message := &Message{}
		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return messages, nil
}

// GetByID retrieves a message by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*Message, error) {
	message := &Message{}
	
	query := `
		SELECT id, user_id, content, created_at, updated_at
		FROM messages
		WHERE id = $1
	`
	
	err := r.db.QueryRow(ctx, query, id).Scan(
		&message.ID,
		&message.UserID,
		&message.Content,
		&message.CreatedAt,
		&message.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}
	
	return message, nil
}

// Update updates a message
func (r *Repository) Update(ctx context.Context, id, userID, content string) error {
	// Check if message exists and belongs to the user
	var ownerID string
	err := r.db.QueryRow(ctx, 
		"SELECT user_id FROM messages WHERE id = $1", 
		id,
	).Scan(&ownerID)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}
	
	// Check ownership
	if ownerID != userID {
		return ErrForbidden
	}
	
	// Update message
	query := `
		UPDATE messages
		SET content = $1, updated_at = NOW()
		WHERE id = $2
	`
	
	_, err = r.db.Exec(ctx, query, content, id)
	return err
}

// Delete deletes a message
func (r *Repository) Delete(ctx context.Context, id, userID string) error {
	// Check if message exists and belongs to the user
	var ownerID string
	err := r.db.QueryRow(ctx, 
		"SELECT user_id FROM messages WHERE id = $1", 
		id,
	).Scan(&ownerID)
	
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMessageNotFound
		}
		return err
	}
	
	// Check ownership
	if ownerID != userID {
		return ErrForbidden
	}
	
	// Delete message
	query := "DELETE FROM messages WHERE id = $1"
	_, err = r.db.Exec(ctx, query, id)
	return err
}