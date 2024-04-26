package handler_test

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"

	"github.com/michelm117/cycling-coach-lab/handler"
	"github.com/michelm117/cycling-coach-lab/mocks"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

func createLoginRequest(email, password string) *http.Request {
	form := url.Values{}
	form.Add("email", email)
	form.Add("password", password)
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestLogin(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userID := 1
	email := "example@example.com"
	password := "password"
	hashedPassword := "1234567890"
	sessionUUID := "281ff6ed-035b-4a9a-a542-fca163d7cde2"

	logger := zaptest.NewLogger(t).Sugar()

	ctrl := gomock.NewController(t)

	mu := mocks.NewMockUserServicer(ctrl)
	mc := mocks.NewMockCryptoer(ctrl)
	ms := mocks.NewMockSessionServicer(ctrl)
	mb := mocks.NewMockBrowserSessionManager(ctrl)
	t.Run("EmailDoesNotExists", func(t *testing.T) {
		mu.EXPECT().GetByEmail(gomock.Eq(email)).Return(nil, sql.ErrNoRows)
		handler := handler.NewAuthHandler(mu, nil, nil, nil, nil, logger)

		// Create a request
		req := createLoginRequest(email, password)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Errorf(t, handler.Login(c), "Invalid credentials")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "")
	})

	t.Run("InactiveUser", func(t *testing.T) {
		mu.EXPECT().GetByEmail(gomock.Eq(email)).Return(&model.User{ID: userID, Status: "inactive"}, nil)
		handler := handler.NewAuthHandler(mu, nil, nil, nil, nil, logger)

		// Create a request
		req := createLoginRequest(email, password)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Errorf(t, handler.Login(c), "Invalid credentials")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "")
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mu.EXPECT().GetByEmail(gomock.Eq(email)).Return(&model.User{ID: userID, Status: "active", PasswordHash: hashedPassword}, nil)
		mc.EXPECT().CompareHashAndPassword(gomock.Eq([]byte(hashedPassword)), gomock.Eq([]byte(password))).Return(fmt.Errorf("Invalid password"))
		handler := handler.NewAuthHandler(mu, nil, nil, nil, mc, logger)

		// Create a request
		req := createLoginRequest(email, password)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Errorf(t, handler.Login(c), "Invalid credentials")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "")
	})

	t.Run("SaveSessionError", func(t *testing.T) {
		mu.EXPECT().GetByEmail(gomock.Eq(email)).Return(&model.User{ID: userID, Status: "active", PasswordHash: hashedPassword}, nil)
		mc.EXPECT().CompareHashAndPassword(gomock.Eq([]byte(hashedPassword)), gomock.Eq([]byte(password))).Return(nil)
		ms.EXPECT().SaveSession(gomock.Eq(userID)).Return("", fmt.Errorf("Error while trying to save session"))
		handler := handler.NewAuthHandler(mu, ms, nil, nil, mc, logger)

		// Create a request
		req := createLoginRequest(email, password)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Errorf(t, handler.Login(c), "Invalid credentials")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "")
	})

	t.Run("SaveBrowserSessionError", func(t *testing.T) {
		mu.EXPECT().GetByEmail(gomock.Eq(email)).Return(&model.User{ID: userID, Status: "active", PasswordHash: hashedPassword}, nil)
		mc.EXPECT().CompareHashAndPassword(gomock.Eq([]byte(hashedPassword)), gomock.Eq([]byte(password))).Return(nil)
		ms.EXPECT().SaveSession(gomock.Eq(userID)).Return(sessionUUID, nil)
		mb.EXPECT().SaveSession(gomock.Any(), gomock.Eq(sessionUUID)).Return(assert.AnError)
		handler := handler.NewAuthHandler(mu, ms, nil, mb, mc, logger)

		// Create a request
		req := createLoginRequest(email, password)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Errorf(t, handler.Login(c), "Invalid credentials")
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "")
	})

	t.Run("ValidCredentials", func(t *testing.T) {
		mu.EXPECT().GetByEmail(gomock.All()).Return(&model.User{ID: userID, Status: "active", PasswordHash: hashedPassword, Email: email}, nil)
		mc.EXPECT().CompareHashAndPassword(gomock.Eq([]byte(hashedPassword)), gomock.Eq([]byte(password))).Return(nil)
		ms.EXPECT().SaveSession(gomock.Eq(userID)).Return(sessionUUID, nil)
		mb.EXPECT().SaveSession(gomock.Any(), gomock.Eq(sessionUUID)).Return(nil)
		handler := handler.NewAuthHandler(mu, ms, nil, mb, mc, logger)

		// Create a request
		req := createLoginRequest(email, password)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.Login(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestLogout(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sessionUUID := "281ff6ed-035b-4a9a-a542-fca163d7cde2"
	logger := zaptest.NewLogger(t).Sugar()
	ctrl := gomock.NewController(t)

	ms := mocks.NewMockSessionServicer(ctrl)
	mb := mocks.NewMockBrowserSessionManager(ctrl)
	t.Run("BrowserSessionDeleteError", func(t *testing.T) {
		mb.EXPECT().DeleteSession(gomock.Any()).Return("", assert.AnError)
		handler := handler.NewAuthHandler(nil, nil, nil, mb, nil, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Errorf(t, handler.Logout(c), "Internal server error")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("SessionDeleteError", func(t *testing.T) {
		mb.EXPECT().DeleteSession(gomock.Any()).Return(sessionUUID, nil)
		ms.EXPECT().DeleteSession(gomock.Eq(sessionUUID)).Return(fmt.Errorf("Error while trying to delete session"))
		handler := handler.NewAuthHandler(nil, ms, nil, mb, nil, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Errorf(t, handler.Logout(c), "Internal server error")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mb.EXPECT().DeleteSession(gomock.Any()).Return(sessionUUID, nil)
		ms.EXPECT().DeleteSession(gomock.Eq(sessionUUID)).Return(nil)
		handler := handler.NewAuthHandler(nil, ms, nil, mb, nil, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.Logout(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestRenderLogin(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	ctrl := gomock.NewController(t)
	mg := mocks.NewMockGlobalSettingServicer(ctrl)

	t.Run("AppNotInitialized", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(false)

		handler := handler.NewAuthHandler(nil, nil, mg, nil, nil, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.RenderLogin(c))
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	})

	t.Run("AppInitialized", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(true)

		handler := handler.NewAuthHandler(nil, nil, mg, nil, nil, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/login", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.RenderLogin(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})
}
