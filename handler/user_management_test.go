package handler_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/michelm117/cycling-coach-lab/handler"
	"github.com/michelm117/cycling-coach-lab/mocks"
	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/test_utils"
)

func TestRenderUserManagementPage(t *testing.T) {
	au := model.User{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com"}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServicer(ctrl)
	mv := mocks.NewMockValidator(ctrl)

	t.Run("Error", func(t *testing.T) {
		us.EXPECT().GetAllUsers().Return(nil, assert.AnError)

		handler := handler.NewUserManagementHandler(us, mv, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

		// Call the handler
		assert.ErrorContains(t, handler.RenderUserManagementPage(c), "Could not retrieve users")
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		users := []*model.User{
			{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com"},
			{ID: 2, Firstname: "Jane", Lastname: "Doe", Email: "jane@doe.com"},
		}
		us.EXPECT().GetAllUsers().Return(users, nil)

		handler := handler.NewUserManagementHandler(us, mv, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

		// Call the handler
		assert.NoError(t, handler.RenderUserManagementPage(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})
}

func TestRenderUserManagementView(t *testing.T) {
	au := model.User{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com"}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	us := mocks.NewMockUserServicer(ctrl)

	t.Run("Error", func(t *testing.T) {
		us.EXPECT().GetAllUsers().Return(nil, assert.AnError)

		handler := handler.NewUserManagementHandler(us, nil, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/users/view", nil)
		rec := httptest.NewRecorder()
		c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

		// Call the handler
		assert.ErrorContains(t, handler.RenderUserManagementView(c), "Could not retrieve users")
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		users := []*model.User{
			{ID: 1, Firstname: "John", Lastname: "Doe", Email: "john@doe.com"},
			{ID: 2, Firstname: "Jane", Lastname: "Doe", Email: "jane@doe.com"},
		}
		us.EXPECT().GetAllUsers().Return(users, nil)

		handler := handler.NewUserManagementHandler(us, nil, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodGet, "/users/view", nil)
		rec := httptest.NewRecorder()
		c := model.AuthenticatedContext{Context: echo.New().NewContext(req, rec), User: &au}

		// Call the handler
		assert.NoError(t, handler.RenderUserManagementView(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServicer(ctrl)

	t.Run("Invalid user id", func(t *testing.T) {
		handler := handler.NewUserManagementHandler(us, nil, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodDelete, "/users/invalid", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		// Call the handler
		assert.ErrorContains(t, handler.DeleteUser(c), "Invalid user id")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Error deleting user", func(t *testing.T) {
		id := 1
		us.EXPECT().DeleteUser(id).Return(assert.AnError)

		handler := handler.NewUserManagementHandler(us, nil, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))

		// Call the handler
		assert.ErrorContains(t, handler.DeleteUser(c), "Could not delete user with id '1'")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("Success", func(t *testing.T) {
		id := 1
		us.EXPECT().DeleteUser(id).Return(nil)

		handler := handler.NewUserManagementHandler(us, nil, nil)

		// Create a request
		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(id))

		// Call the handler
		assert.NoError(t, handler.DeleteUser(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func createAddUserRequest(firstname, lastname, email, password, confirmPassword, dateOfBirth, role string) *http.Request {
	form := url.Values{}
	form.Add("firstname", firstname)
	form.Add("lastname", lastname)
	form.Add("email", email)
	form.Add("password", password)
	form.Add("confirmPassword", confirmPassword)
	form.Add("dateOfBirth", dateOfBirth)
	form.Add("role", role)

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func TestRenderAddUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	us := mocks.NewMockUserServicer(ctrl)
	mv := mocks.NewMockValidator(ctrl)

	t.Run("error invalid first name", func(t *testing.T) {
		mv.EXPECT().CreateValidUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
		handler := handler.NewUserManagementHandler(us, mv, nil)

		// Create a request
		req := createAddUserRequest("", "Doe", "john@doe.com", "password", "password", "1990-01-01", "admin")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.Error(t, handler.AddUser(c))
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("error adding user", func(t *testing.T) {
		us.EXPECT().AddUser(gomock.Any()).Return(nil, assert.AnError)
		mv.EXPECT().CreateValidUser(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(&model.User{}, nil)
		handler := handler.NewUserManagementHandler(us, mv, nil)

		// Create a request
		req := createAddUserRequest("John", "Doe", "john@doe.com", "password", "password", "1990-01-01", "admin")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.ErrorContains(t, handler.AddUser(c), "Could not add user")
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("success", func(t *testing.T) {
		user := &model.User{
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john@doe.com",
		}
		us.EXPECT().AddUser(gomock.Any()).Return(&model.User{ID: 1}, nil)
		mv.EXPECT().CreateValidUser(gomock.Eq("John"), gomock.Eq("Doe"), gomock.Eq("admin"), gomock.Eq("john@doe.com"), gomock.Eq("1990-01-01"), gomock.Eq("password"), gomock.Eq("password")).Return(user, nil)

		handler := handler.NewUserManagementHandler(us, mv, nil)

		// Create a request
		req := createAddUserRequest("John", "Doe", "john@doe.com", "password", "password", "1990-01-01", "admin")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		// Call the handler
		assert.NoError(t, handler.AddUser(c))
		assert.Equal(t, http.StatusOK, rec.Code)
		test_utils.MakeSnapshot(t, rec.Body.String())
	})
}
