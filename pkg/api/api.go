package api

import (
	"chatapp/internal/config"
	"chatapp/internal/service/email"
	"context"
)

func Run(ctx context.Context) error {
	cfg := config.MustInit()

	emailService := email.NewEmailServiceSMTP(cfg.Email)

	emailService.SendMagicLink(ctx, "test@test.test", "https://google.com")

	return nil
}
