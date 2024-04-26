package services_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/assert"

	"github.com/michelm117/cycling-coach-lab/services"
)

func SetUpTestSession(userID int, sessionService services.SessionServicer, t *testing.T) string {
	sessionID, err := sessionService.SaveSession(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, sessionID)
	return sessionID
}

func TestSaveAndGetByUUID(t *testing.T) {
	userID := 1
	sessionService := services.NewSessionServicer(DB, nil)
	sessionID := SetUpTestSession(userID, sessionService, t)

	session, err := sessionService.GetByUUID(sessionID)

	assert.NoError(t, err)
	assert.NotNil(t, session)
	assert.Equal(t, userID, session.UserID)
}

func TestDeleteSession(t *testing.T) {
	userID := 1
	sessionService := services.NewSessionServicer(DB, nil)
	sessionID := SetUpTestSession(userID, sessionService, t)

	err := sessionService.DeleteSession(sessionID)
	assert.NoError(t, err)

	// Verify session is deleted
	_, err = sessionService.GetByUUID(sessionID)
	assert.Error(t, err) // Expecting an error as session should not exist
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestGetByUUIDWhenSessionDoesNotExist(t *testing.T) {
	sessionService := services.NewSessionServicer(DB, nil)
	uuid := uuid.New().String()
	_, err := sessionService.GetByUUID(uuid)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestDeleteExpiredSessions(t *testing.T) {
	// manually add expired session
	sessionID := uuid.New().String()
	createdAt := time.Date(2022, time.January, 15, 12, 0, 0, 0, time.UTC)

	stmt, err := DB.Prepare("INSERT INTO sessions (id, created_at, user_id) VALUES ($1, $2, $3)")
	defer stmt.Close()
	assert.NoError(t, err)

	_, err = stmt.Exec(sessionID, createdAt, "1")
	assert.NoError(t, err)

	sessionService := services.NewSessionServicer(DB, nil)
	beforeUserSessions, err := sessionService.GetByUserID(1)

	err = sessionService.DeleteExpiredSessions()
	assert.NoError(t, err)

	afterUserSessions, err := sessionService.GetByUserID(1)
	assert.NoError(t, err)

	assert.Greater(t, len(beforeUserSessions), len(afterUserSessions))
}

func TestGetByUserID(t *testing.T) {
	sessionService := services.NewSessionServicer(DB, nil)
	userID := 1
	sessionID := SetUpTestSession(userID, sessionService, t)

	sessions, err := sessionService.GetByUserID(userID)
	assert.NoError(t, err)
	assert.Greater(t, len(sessions), 0)

	found := false
	for _, s := range sessions {
		if s.ID == sessionID {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestAutenticateUserBySessionID_WhenSessionExists(t *testing.T) {
	sessionService := services.NewSessionServicer(DB, nil)
	userService := services.NewUserServicer(DB, nil)
	user, err := userService.GetByEmail("admin@example.com")
	assert.NoError(t, err)

	sessionID, err := sessionService.SaveSession(user.ID)
	assert.NoError(t, err)

	authenticatedUser, err := sessionService.AuthenticateUserBySessionID(sessionID)
	assert.NoError(t, err)
	assert.NotNil(t, authenticatedUser)
	assert.Equal(t, user.ID, authenticatedUser.ID)
}

func TestAutenticateUserBySessionID_WhenSessionDoesNotExist(t *testing.T) {
	sessionService := services.NewSessionServicer(DB, nil)
	_, err := sessionService.AuthenticateUserBySessionID(uuid.New().String())
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestAutenticateUserBySessionID_WhenEmptySessionID(t *testing.T) {
	sessionService := services.NewSessionServicer(DB, nil)
	_, err := sessionService.AuthenticateUserBySessionID("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session ID is empty")
}

func TestAutenticateUserBySessionID_WhenSessionIsExpired(t *testing.T) {
	sessionService := services.NewSessionServicer(DB, nil)
	sessionID := uuid.New().String()
	createdAt := time.Date(2022, time.January, 15, 12, 0, 0, 0, time.UTC)

	stmt, err := DB.Prepare("INSERT INTO sessions (id, created_at, user_id) VALUES ($1, $2, $3)")
	defer stmt.Close()
	assert.NoError(t, err)

	_, err = stmt.Exec(sessionID, createdAt, "1")
	assert.NoError(t, err)

	_, err = sessionService.AuthenticateUserBySessionID(sessionID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}
