package services

import (
	"database/sql"
	"fmt"

	gomail "github.com/Shopify/gomail"
	"go.uber.org/zap"

	"github.com/michelm117/cycling-coach-lab/model"
	"github.com/michelm117/cycling-coach-lab/utils"
)

type EmailServicer interface {
	GetEmailSettings() (*model.EmailSettings, error)
	SaveEmailSettings(*model.EmailSettings) error
	SendEmail(to []string, subject, body string) error
	SendEmailWithTemplate(to []string, subject, templatePath string, data interface{}) error
}

type EmailService struct {
	globalSettingServicer GlobalSettingServicer
	templater             Templater
	logger                *zap.SugaredLogger
}

func NewEmailService(globalSettingServicer GlobalSettingServicer, templater Templater, logger *zap.SugaredLogger,
) EmailServicer {
	return &EmailService{
		globalSettingServicer: globalSettingServicer,
		templater:             templater,
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

	m := gomail.NewMessage()
	m.SetHeader("From", settings.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	emailSender, err := utils.NewGomailSender(settings.Host, settings.Port, settings.Username, settings.Password)
	if err != nil {
		return fmt.Errorf("Invalid email settings")
	}

	if err := emailSender.Send(settings.From, to, subject, body); err != nil {
		s.logger.Errorw("Failed to send email", "error", err)
		return err
	}
	return nil
}

func (s *EmailService) SendEmailWithTemplate(to []string, subject, templatePath string, data interface{}) error {
	body, err := s.templater.Template(templatePath, data)
	if err != nil {
		return err
	}
	return s.SendEmail(to, subject, body)
}
