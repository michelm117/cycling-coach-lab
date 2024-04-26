package utils

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

const (
	INFO    = "info"
	SUCCESS = "success"
	WARNING = "warning"
	DANGER  = "danger"
)

type Toast struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func New(level string, message string) Toast {
	return Toast{level, message}
}

func Info(message string) Toast {
	return New(INFO, message)
}

func Success(c echo.Context, message string) error {
	New(SUCCESS, message).SetHXTriggerHeader(c)
	return nil
}

func Warning(message string) Toast {
	return New(WARNING, message)
}

func Danger(message string) Toast {
	return New(DANGER, message)
}

func (t Toast) Error() string {
	return fmt.Sprintf("%s", t.Message)
}

func (t Toast) jsonify() (string, error) {
	t.Message = t.Error()
	eventMap := map[string]Toast{}
	eventMap["showToast"] = t
	jsonData, err := json.Marshal(eventMap)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (t Toast) SetHXTriggerHeader(c echo.Context) {
	jsonData, _ := t.jsonify()
	c.Response().Header().Set("HX-Trigger", jsonData)
}
