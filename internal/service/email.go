package service

import (
	"chatapp/internal/config"
	"context"
	"net/smtp"
)

type EmailService interface {
	SendMagicLink(ctx context.Context, email, magicLink string) error
}

type EmailServiceSMTP struct {
	host string
	port int
	auth smtp.Auth
	from string
}

func NewEmailServiceSMTP(cfg config.EmailConfig) EmailService {
	return &EmailServiceSMTP{
		host: cfg.Host,
		port: cfg.Port,
		from: cfg.From,
		auth: smtp.PlainAuth("", cfg.User, cfg.Pass, cfg.Host),
	}
}

// SendMagicLink implements [EmailService].
func (e *EmailServiceSMTP) SendMagicLink(ctx context.Context, to string, magicLink string) error {
	return nil
}
