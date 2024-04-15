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

type GlobalSettingService struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func NewGlobalSettingService(db *sql.DB, logger *zap.SugaredLogger) *GlobalSettingService {
	return &GlobalSettingService{
		db:     db,
		logger: logger,
	}
}

func (s *GlobalSettingService) Create(setting *model.GlobalSetting) error {
	query := `
		INSERT INTO Configuration.GlobalSettings(SectionName, SettingName, SettingValue, SettingType)
		VALUES (?, ?, ?, ?)
	`
	_, err := s.db.Exec(query,
		strings.ToLower(setting.SectionName),
		strings.ToLower(setting.SettingName),
		setting.SettingValue,
		setting.SettingType)
	if err != nil {
		s.logger.Errorw("Failed to create global setting", "error", err)
		return fmt.Errorf("failed to create global setting: %w", err)
	}
	return nil
}

func (s *GlobalSettingService) GetBySectionAndName(sectionName, settingName string) (interface{}, error) {
	query := `
		SELECT SectionName, SettingName, SettingValue, SettingType
		FROM Configuration.GlobalSettings
		WHERE SectionName = ? AND SettingName = ?
	`
	row := s.db.QueryRow(query, sectionName, settingName)
	setting := &model.GlobalSetting{}
	err := row.Scan(&setting.SectionName, &setting.SettingName, &setting.SettingValue, &setting.SettingType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("global setting with section '%s' and name '%s' not found", sectionName, settingName)
		}
		s.logger.Errorw("Failed to get global setting", "error", err)
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
			s.logger.Errorf("Failed to parse setting %s of type %d", settingValue, settingType)
			return nil, fmt.Errorf("failed to parse setting %s of type %d", settingValue, settingType)
		}
		return value, nil
	case booleanSetting:
		value, err := strconv.ParseBool(settingValue)
		if err != nil {
			s.logger.Errorf("Failed to parse setting %s of type %d", settingValue, settingType)
			return nil, fmt.Errorf("failed to parse setting %s of type %d", settingValue, settingType)
		}
		return value, nil
	default:
		s.logger.Errorf("Unknown setting type '%d' for setting %s", settingType, settingValue)
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
