package utils

import (
	"crypto/tls"
	"strconv"

	"github.com/Shopify/gomail"
)

// EmailSettings holds the settings for the email service
type EmailSettings struct {
	From     string
	Username string
	Password string
	Host     string
	Port     string
}

// Logger is a placeholder for your logging mechanism
type Logger interface {
	Errorw(message string, keysAndValues ...interface{})
}

// EmailSender is an interface for sending emails
type EmailSender interface {
	Send(from string, to []string, subject, body string) error
}

// GomailSender is an implementation of EmailSender using gomail
type GomailSender struct {
	host     string
	port     int
	username string
	password string
}

// Send sends an email using gomail
func (g *GomailSender) Send(from string, to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer(g.host, g.port, g.username, g.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return d.DialAndSend(m)
}

// NewGomailSender creates a new GomailSender
func NewGomailSender(host, port, username, password string) (*GomailSender, error) {
	p, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}

	return &GomailSender{
		host:     host,
		port:     p,
		username: username,
		password: password,
	}, nil
}
