package services

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
)

type SessionService struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewSessionService(db *sql.DB, logger *zap.SugaredLogger) *SessionService {
	return &SessionService{
		db:     db,
		logger: logger,
	}
}

func (s *SessionService) GetByUUID(uuid string) (*model.Session, error) {
	row := s.db.QueryRow("SELECT * FROM sessions WHERE sessions.id = $1", uuid)

	var session model.Session
	err := row.Scan(
		&session.ID,
		&session.UserID,
		&session.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session with id '%s' not found", uuid)
		}
		return nil, fmt.Errorf("error while trying to get session with uuid '%s':\n%s", uuid, err)
	}
	return &session, nil
}

func (s *SessionService) SaveSession(userID int) (string, error) {
	sessionID := uuid.New().String()
	stmt, err := s.db.Prepare("INSERT INTO sessions (id, user_id) VALUES ($1, $2)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sessionID, userID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *SessionService) GetUserBySession(sessionID string) (*model.User, error) {
	query := `
        SELECT u.*
        FROM users u
        INNER JOIN sessions s ON u.id = s.user_id
        WHERE s.id = $1
    `
	row := s.db.QueryRow(query, sessionID)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Firstname,
		&user.Lastname,
		&user.DateOfBirth,
		&user.PasswordHash,
		&user.Status,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with session id '%s' not found", sessionID)
		}
		return nil, fmt.Errorf("error while trying to execute query: %s\n%s", query, err)
	}

	return &user, nil
}
