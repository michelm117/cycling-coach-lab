package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
)

type SessionServicer interface {
	GetByUUID(uuid string) (*model.Session, error)
	GetByUserID(userID int) ([]model.Session, error)
	AuthenticateUserBySessionID(sessionID string) (*model.User, error)
	SaveSession(userID int) (string, error)
	DeleteSession(sessionID string) error
	DeleteExpiredSessions() error
	DeleteAllSessions() error
	ScheduleSessionCleanup()
}

type SessionService struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewSessionServicer(db *sql.DB, logger *zap.SugaredLogger) SessionServicer {
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
		&session.CreatedAt,
		&session.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetUserBySession retrieves a user associated with the given session ID.
// It performs a database query to fetch the user details based on the session ID.
// The query joins the 'users' and 'sessions' tables and filters users by session ID,
// active status, and session creation date not older than the specified interval in days.
//
// Parameters:
//   - sessionID: The session ID used to identify the session.
//
// Returns:
//   - *model.User: A pointer to the user object if found.
//   - error: An error if any occurred during the database operation.
func (s *SessionService) AuthenticateUserBySessionID(sessionID string) (*model.User, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID is empty")
	}

	query := `
        SELECT u.*
        FROM users u
        INNER JOIN sessions s ON u.id = s.user_id
        WHERE s.id = $1
          AND s.created_at >= NOW() - INTERVAL '7 days' 
          AND u.status = 'active';
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
		return nil, err
	}

	return &user, nil
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

func (s *SessionService) DeleteSession(sessionID string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE id = $1", sessionID)
	if err != nil {
		return fmt.Errorf("error while trying to delete session with id '%s':\n%s", sessionID, err)
	}
	return nil
}

// DeleteExpiredSessions deletes old expired sessions from the database.
func (s *SessionService) DeleteExpiredSessions() error {
	query := `DELETE FROM sessions WHERE created_at < NOW() - INTERVAL '7 days';`
	_, err := s.db.Exec(query)
	if err != nil {
		return fmt.Errorf("could not delete expired sessions: %w", err)
	}
	return nil
}

// ScheduleSessionCleanup schedules a new ticker that triggers every 6 hours to clean up expired sessions.
func (s *SessionService) ScheduleSessionCleanup() {
	go func() {
		ticker := time.NewTicker(time.Hour * 6)
		defer ticker.Stop()

		for {
			<-ticker.C
			if err := s.DeleteExpiredSessions(); err != nil {
				s.logger.Error("error cleaning up sessions", err)
			}
		}
	}()
}

func (s *SessionService) GetByUserID(userID int) ([]model.Session, error) {
	rows, err := s.db.Query("SELECT * FROM sessions WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.Session
	for rows.Next() {
		var session model.Session
		err := rows.Scan(
			&session.ID,
			&session.CreatedAt,
			&session.UserID,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

func (s *SessionService) DeleteAllSessions() error {
	_, err := s.db.Exec("DELETE FROM sessions")
	if err != nil {
		return fmt.Errorf("error while trying to delete all sessions:\n%s", err)
	}
	return nil
}
