package services_test

import (
	"fmt"
	"testing"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/stretchr/testify/assert"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/services"
)

func TestCreate_ValidInput(t *testing.T) {
	settingService := services.NewGlobalSettingService(DB, nil)
	err := settingService.Create(&model.GlobalSetting{
		SectionName:  "section",
		SettingName:  "name",
		SettingValue: "value",
		SettingType:  0,
	})

	assert.Nil(t, err)

	// Check if the setting was created
	setting, err := settingService.GetBySectionAndName("section", "name")
	assert.Nil(t, err)
	assert.NotNil(t, setting)
	assert.Equal(t, "value", setting)
}

func TestCreate_InvalidInput(t *testing.T) {
	settingService := services.NewGlobalSettingService(DB, nil)
	err := settingService.Create(&model.GlobalSetting{
		SectionName:  "section",
		SettingName:  "name",
		SettingValue: "value",
		SettingType:  -10,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid setting type: -10")
}

func TestStringGetBySectionAndName(t *testing.T) {
	testCases := []struct {
		name          string
		settingType   int
		settingValue  string
		expectedValue interface{}
	}{
		{
			name:          "string",
			settingType:   0,
			settingValue:  "value",
			expectedValue: "value",
		},
		{
			name:          "integer",
			settingType:   1,
			settingValue:  "10",
			expectedValue: 10,
		},
		{
			name:          "boolean",
			settingType:   2,
			settingValue:  "true",
			expectedValue: true,
		},
	}

	settingService := services.NewGlobalSettingService(DB, nil)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestGetBySectionAndName_%s", tc.name), func(t *testing.T) {
			err := settingService.Create(&model.GlobalSetting{
				SectionName:  "section",
				SettingName:  fmt.Sprintf("name_%s", tc.name),
				SettingValue: tc.settingValue,
				SettingType:  tc.settingType,
			})
			assert.Nil(t, err)

			setting, err := settingService.GetBySectionAndName("section", fmt.Sprintf("name_%s", tc.name))
			assert.Nil(t, err)
			assert.NotNil(t, setting)
			assert.Equal(t, tc.expectedValue, setting)
		})
	}
}

func TestParseSettingsValue(t *testing.T) {
	testCases := []struct {
		name          string
		settingType   int
		settingValue  string
		expectedValue interface{}
	}{
		{
			name:          "string",
			settingType:   0,
			settingValue:  "value",
			expectedValue: "value",
		},
		{
			name:          "integer",
			settingType:   1,
			settingValue:  "10",
			expectedValue: 10,
		},
		{
			name:          "boolean",
			settingType:   2,
			settingValue:  "true",
			expectedValue: true,
		},
	}

	settingService := services.NewGlobalSettingService(DB, nil)
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("TestParseSettingsValue_%s", tc.name), func(t *testing.T) {
			value, err := settingService.ParseSettingsValue(tc.settingType, tc.settingValue)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedValue, value)
		})
	}
}

func TestIsAppInitialzed(t *testing.T) {
	settingService := services.NewGlobalSettingService(DB, nil)

	initialized := settingService.IsAppInitialized()
	assert.False(t, initialized)

	err := settingService.InitializeApp()
	assert.Nil(t, err)
	initialized = settingService.IsAppInitialized()
	assert.True(t, initialized)
}

func TestLowerCaseNamesAreSaved(t *testing.T) {
	settingService := services.NewGlobalSettingService(DB, nil)
	t.Run("Upper case section & setting name is save lowercase", func(t *testing.T) {
		err := settingService.Create(&model.GlobalSetting{
			SectionName:  "SECTION_NAME",
			SettingName:  "SETTING_NAME",
			SettingValue: "true",
			SettingType:  2,
		})
		assert.Nil(t, err)

		val, err := settingService.GetBySectionAndName("SECTION_NAME", "SETTING_NAME")
		assert.Nil(t, err)
		assert.True(t, val.(bool))
	})
}

func TestUpdate(t *testing.T) {
	settingService := services.NewGlobalSettingService(DB, nil)
	err := settingService.Create(&model.GlobalSetting{
		SectionName:  "section",
		SettingName:  "name",
		SettingValue: "value",
		SettingType:  0,
	})
	assert.Nil(t, err)

	err = settingService.Create(&model.GlobalSetting{
		SectionName:  "section",
		SettingName:  "name",
		SettingValue: "new_value",
		SettingType:  0,
	})
	assert.Nil(t, err)

	val, err := settingService.GetBySectionAndName("section", "name")
	assert.Nil(t, err)
	assert.Equal(t, "new_value", val)
}
