package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"

	"github.com/michelm117/cycling-coach-lab/handler"
	"github.com/michelm117/cycling-coach-lab/mocks"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

func TestRenderSetup(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mg := mocks.NewMockGlobalSettingServicer(ctrl)

	t.Run("App already initialized", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(true)

		handler := handler.NewSetupHandler(mg, nil, nil, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/setup", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.RenderSetup(c))
		assert.Equal(t, http.StatusTemporaryRedirect, rec.Code)
	})

	t.Run("App not initialized", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(false)

		handler := handler.NewSetupHandler(mg, nil, nil, logger)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/setup", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.RenderSetup(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})
}

func createSetupRequest(firstname, lastname, email, password, dateOfBirth string) *http.Request {
	form := url.Values{}
	form.Add("firstname", firstname)
	form.Add("lastname", lastname)
	form.Add("email", email)
	form.Add("password", password)
	form.Add("dateOfBirth", dateOfBirth)

	req := httptest.NewRequest(http.MethodPost, "/setup", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestSetup(t *testing.T) {
	logger := zaptest.NewLogger(t).Sugar()
	ctrl := gomock.NewController(t)
	mg := mocks.NewMockGlobalSettingServicer(ctrl)
	mc := mocks.NewMockCryptoer(ctrl)
	mu := mocks.NewMockUserServicer(ctrl)

	t.Run("App already initialized", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(true)

		handler := handler.NewSetupHandler(mg, nil, nil, logger)

		// Create a request
		req := createSetupRequest("John", "Doe", "john@doe.com", "password", "1990-01-01")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.ErrorContains(t, handler.Setup(c), "App already initialized")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Invalid date of birth", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(false)

		handler := handler.NewSetupHandler(mg, nil, nil, logger)

		// Create a request
		req := createSetupRequest("John", "Doe", "john@doe.com", "password", "invalid")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.ErrorContains(t, handler.Setup(c), "Invalid date of birth")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Internal server error on cryptoer generate hash", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(false)
		mc.EXPECT().GenerateFromPassword(gomock.Eq([]byte("password"))).Return(nil, assert.AnError)

		handler := handler.NewSetupHandler(mg, nil, mc, logger)

		// Create a request
		req := createSetupRequest("John", "Doe", "john@doe.com", "password", "1990-01-01")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.ErrorContains(t, handler.Setup(c), "Internal server error")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Internal server error on db save", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(false)
		mc.EXPECT().GenerateFromPassword(gomock.Eq([]byte("password"))).Return([]byte("hashed"), nil)
		mu.EXPECT().AddUser(gomock.Any()).Return(nil, assert.AnError)

		handler := handler.NewSetupHandler(mg, mu, mc, logger)

		// Create a request
		req := createSetupRequest("John", "Doe", "john@doe.com", "password", "1990-01-01")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.ErrorContains(t, handler.Setup(c), "Internal server error")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Internal server error on initialising app", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(false)
		mg.EXPECT().InitializeApp().Return(assert.AnError)
		mc.EXPECT().GenerateFromPassword(gomock.Eq([]byte("password"))).Return([]byte("hashed"), nil)
		mu.EXPECT().AddUser(gomock.Any()).Return(nil, nil)

		handler := handler.NewSetupHandler(mg, mu, mc, logger)

		// Create a request
		req := createSetupRequest("John", "Doe", "john@doe.com", "password", "1990-01-01")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.ErrorContains(t, handler.Setup(c), "Internal server error")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Success", func(t *testing.T) {
		mg.EXPECT().IsAppInitialized().Return(false)
		mg.EXPECT().InitializeApp().Return(nil)
		mc.EXPECT().GenerateFromPassword(gomock.Eq([]byte("password"))).Return([]byte("hashed"), nil)
		mu.EXPECT().AddUser(gomock.Any()).Return(nil, nil)

		handler := handler.NewSetupHandler(mg, mu, mc, logger)

		// Create a request
		req := createSetupRequest("John", "Doe", "john@doe.com", "password", "1990-01-01")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.Setup(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
