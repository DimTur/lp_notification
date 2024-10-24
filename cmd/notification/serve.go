package notification

import (
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DimTur/lp_notification/internal/app/telegram"
	"github.com/DimTur/lp_notification/internal/config"
	rabbitmq_store "github.com/DimTur/lp_notification/internal/storage/rabbitmq"
	"github.com/spf13/cobra"
)

func NewServeCmd() *cobra.Command {
	var configPath string

	c := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"s"},
		Short:   "Start API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

			ctx, cancel := signal.NotifyContext(cmd.Context(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
			defer cancel()
			var wg sync.WaitGroup

			cfg, err := config.Parse(configPath)
			if err != nil {
				return err
			}

			storage, err := rabbitmq_store.NewRabbit()
			if err != nil {
				return err
			}

			// start tg bot
			wg.Add(1)
			go func() {
				defer wg.Done()
				telegram.RunTg(
					ctx,
					cfg.TelegramBot.TgBotHost,
					cfg.TelegramBot.TgBotToken,
					cfg.TelegramBot.BatchSize,
					storage,
					log,
				)
			}()

			log.Info("tg bot starting at:", slog.Any("port", cfg.Server.Port))
			<-ctx.Done()
			wg.Wait()

			return nil
		}}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
