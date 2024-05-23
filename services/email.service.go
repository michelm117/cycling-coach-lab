package services

import (
	"database/sql"
	"fmt"
	"net/smtp"
	"strings"

	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
)

type EmailServicer interface {
	GetEmailSettings() (*model.EmailSettings, error)
	SaveEmailSettings(*model.EmailSettings) error
	SendEmail(to []string, subject, body string) error
}

type EmailService struct {
	globalSettingServicer GlobalSettingServicer
	logger                *zap.SugaredLogger
}

func NewEmailService(globalSettingServicer GlobalSettingServicer, logger *zap.SugaredLogger) EmailServicer {
	return &EmailService{
		globalSettingServicer: globalSettingServicer,
		logger:                logger,
	}
}

func (s *EmailService) GetEmailSettings() (*model.EmailSettings, error) {
	settings := []string{"from", "username", "password", "host", "port"}
	values := make(map[string]string, len(settings))

	for _, setting := range settings {
		value, err := s.getEmailSetting(setting)
		if err != nil {
			return nil, err
		}
		values[setting] = value
	}

	return &model.EmailSettings{
		From:     values["from"],
		Username: values["username"],
		Password: values["password"],
		Host:     values["host"],
		Port:     values["port"],
	}, nil
}

func (s *EmailService) getEmailSetting(name string) (string, error) {
	value, err := s.globalSettingServicer.GetBySectionAndName("email", name)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		s.logger.Errorw("Failed to get email setting", "name", name, "error", err)
		return "", err
	}
	return value.(string), nil
}

func (s *EmailService) SaveEmailSettings(emailSettings *model.EmailSettings) error {
	settings := map[string]string{
		"from":     emailSettings.From,
		"username": emailSettings.Username,
		"password": emailSettings.Password,
		"host":     emailSettings.Host,
		"port":     emailSettings.Port,
	}

	for setting, value := range settings {
		err := s.globalSettingServicer.Create(&model.GlobalSetting{
			SectionName:  "email",
			SettingName:  setting,
			SettingValue: value,
			SettingType:  stringSetting,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *EmailService) SendEmail(to []string, subject, body string) error {
	settings, err := s.GetEmailSettings()
	if err != nil {
		s.logger.Errorw("Failed to get email settings", "error", err)
		return err
	}

	receiver := ""
	if len(to) > 1 {
		receiver = strings.Join(to, ",")
	} else {
		receiver = to[0]
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", settings.From, receiver, subject, body)
	auth := smtp.PlainAuth("", settings.Username, settings.Password, settings.Host)
	if err := smtp.SendMail(fmt.Sprintf("[%s]:%s", settings.Host, settings.Port), auth, settings.From, to, []byte(msg)); err != nil {
		s.logger.Errorw("Failed to send email", "to", to, "subject", subject, "error", err)
		return err
	}
	return nil
}
