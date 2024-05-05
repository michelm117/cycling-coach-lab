package services

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
)

const (
	stringSetting  = 0
	integerSetting = 1
	booleanSetting = 2
)

type GlobalSettingServicer interface {
	Create(setting *model.GlobalSetting) error
	GetBySectionAndName(sectionName, settingName string) (interface{}, error)
	ParseSettingsValue(settingType int, settingValue string) (interface{}, error)
	IsAppInitialized() bool
	InitializeApp() error
}

type GlobalSettingService struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewGlobalSettingService(db *sql.DB, logger *zap.SugaredLogger) GlobalSettingServicer {
	return &GlobalSettingService{
		db:     db,
		logger: logger,
	}
}

func (s *GlobalSettingService) Create(setting *model.GlobalSetting) error {
	if setting.SettingType < 0 || setting.SettingType > 2 {
		return fmt.Errorf("invalid setting type: %d", setting.SettingType)
	}

	query := `
    INSERT INTO globalSettings(SectionName, SettingName, SettingValue, SettingType)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (SectionName, SettingName) DO UPDATE
        SET SettingValue = EXCLUDED.SettingValue,
            SettingType = EXCLUDED.SettingType
	`
	_, err := s.db.Exec(query,
		setting.SectionName,
		setting.SettingName,
		setting.SettingValue,
		setting.SettingType)
	if err != nil {
		return fmt.Errorf("failed to create/update global setting: %w", err)
	}
	return nil
}

func (s *GlobalSettingService) GetBySectionAndName(sectionName, settingName string) (interface{}, error) {
	query := `
		SELECT SectionName, SettingName, SettingValue, SettingType
		FROM globalSettings
		WHERE SectionName = $1 AND SettingName = $2
	`
	row := s.db.QueryRow(query, strings.ToLower(sectionName), strings.ToLower(settingName))
	setting := &model.GlobalSetting{}
	err := row.Scan(&setting.SectionName, &setting.SettingName, &setting.SettingValue, &setting.SettingType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("global setting with section '%s' and name '%s' not found", sectionName, settingName)
		}
		return nil, fmt.Errorf("failed to get global setting: %w", err)
	}

	value, err := s.ParseSettingsValue(setting.SettingType, setting.SettingValue)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (s *GlobalSettingService) ParseSettingsValue(settingType int, settingValue string) (interface{}, error) {
	switch settingType {
	case stringSetting:
		return settingValue, nil
	case integerSetting:
		value, err := strconv.Atoi(settingValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse setting %s of type %d", settingValue, settingType)
		}
		return value, nil
	case booleanSetting:
		value, err := strconv.ParseBool(settingValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse setting %s of type %d", settingValue, settingType)
		}
		return value, nil
	default:
		return false, fmt.Errorf("Unknown setting type '%d' for setting %s", settingType, settingValue)
	}
}

func (s *GlobalSettingService) IsAppInitialized() bool {
	val, err := s.GetBySectionAndName("app", "initialized")
	if err != nil {
		return false
	}
	return val.(bool)
}

func (s *GlobalSettingService) InitializeApp() error {
	setting := &model.GlobalSetting{
		SectionName:  "app",
		SettingName:  "initialized",
		SettingValue: "true",
		SettingType:  booleanSetting,
	}
	return s.Create(setting)
}
