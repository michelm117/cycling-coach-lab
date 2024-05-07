package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/michelm117/cycling-coach-lab/mocks"
	"github.com/michelm117/cycling-coach-lab/utils"
)

func TestValidateNonEmptyStringField(t *testing.T) {
	v := utils.NewValidator(nil)
	testCases := []struct {
		fieldName string
		value     string
		expected  error
	}{
		{
			fieldName: "test",
			value:     "value",
			expected:  nil,
		},
		{
			fieldName: "test",
			value:     "",
			expected:  fmt.Errorf("Invalid test"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.fieldName, func(t *testing.T) {
			result := v.ValidateNonEmptyStringField(tc.fieldName, tc.value)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateRole(t *testing.T) {
	v := utils.NewValidator(nil)

	testCases := []struct {
		role     string
		expected error
	}{
		{
			role:     "admin",
			expected: nil,
		},
		{
			role:     "athlete",
			expected: nil,
		},
		{
			role:     "invalid",
			expected: fmt.Errorf("Invalid role"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.role, func(t *testing.T) {
			result := v.ValidateRole(tc.role)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidateEmail(t *testing.T) {
	v := utils.NewValidator(nil)

	testCases := []struct {
		email    string
		expected error
	}{
		{
			email:    "valide@email.com",
			expected: nil,
		},
		{
			email:    "invalid",
			expected: fmt.Errorf("Invalid email"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.email, func(t *testing.T) {
			result := v.ValidateEmail(tc.email)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestValidatePassword(t *testing.T) {
	v := utils.NewValidator(nil)

	testCases := []struct {
		password        string
		confirmPassword string
		expected        error
	}{
		{
			password:        "password",
			confirmPassword: "password",
			expected:        nil,
		},
		{
			password:        "password",
			confirmPassword: "different",
			expected:        fmt.Errorf("Passwords do not match"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.password, func(t *testing.T) {
			result := v.ValidatePassword(tc.password, tc.confirmPassword)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCreateValidUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cryptoer := mocks.NewMockCryptoer(ctrl)
	v := utils.NewValidator(cryptoer)
	cryptoer.EXPECT().GenerateFromPassword(gomock.Any()).Return([]byte("hashed"), nil)
	cryptoer.EXPECT().GenerateFromPassword(gomock.Any()).Return([]byte("hashed"), nil)
	testCases := []struct {
		firstname       string
		lastname        string
		role            string
		email           string
		dateOfBirthStr  string
		password        string
		confirmPassword string
		expected        error
	}{
		{
			firstname:       "first",
			lastname:        "last",
			role:            "admin",
			email:           "john@doe.com",
			dateOfBirthStr:  "2000-01-01",
			password:        "password",
			confirmPassword: "password",
			expected:        nil,
		},
		{
			firstname:       "",
			lastname:        "last",
			role:            "admin",
			email:           "john@doe.com",
			dateOfBirthStr:  "2000-01-01",
			password:        "password",
			confirmPassword: "password",
			expected:        fmt.Errorf("Invalid first name"),
		},
		{
			firstname:       "athlete",
			lastname:        "last",
			role:            "athlete",
			email:           "john@doe.com",
			dateOfBirthStr:  "2000-01-01",
			password:        "password",
			confirmPassword: "password",
			expected:        nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.firstname, func(t *testing.T) {
			result, err := v.CreateValidUser(
				tc.firstname,
				tc.lastname,
				tc.role,
				tc.email,
				tc.dateOfBirthStr,
				tc.password,
				tc.confirmPassword,
			)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}
