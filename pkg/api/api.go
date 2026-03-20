package api

import (
	"chatapp/internal/config"
	"chatapp/internal/database"
	"chatapp/internal/service"
	"context"
)

func Run(ctx context.Context) error {
	cfg := config.MustInit()

	db, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	emailService := service.NewEmailServiceSMTP(cfg.Email)
	emailService.SendMagicLink(ctx, "test@test.test", "https://google.com")

	return nil
}
