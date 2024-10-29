package notification

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/DimTur/lp_notification/internal/app/sender"
	"github.com/DimTur/lp_notification/internal/app/telegram"
	tgclient "github.com/DimTur/lp_notification/internal/clients/telegram"
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

			tgClient, err := tgclient.NewTgClient(
				cfg.TelegramBot.TgBotHost,
				cfg.TelegramBot.TgBotToken,
				log,
			)
			if err != nil {
				log.Error("failed init tg client", slog.Any("err", err))
			}

			// Init RabbitMQ
			rmqUrl := fmt.Sprintf(
				"amqp://%s:%s@%s:%d/",
				cfg.RabbitMQ.UserName,
				cfg.RabbitMQ.Password,
				cfg.RabbitMQ.Host,
				cfg.RabbitMQ.Port,
			)
			rmq, err := rabbitmq_store.NewClient(rmqUrl)
			if err != nil {
				log.Error("failed init rabbit mq", slog.Any("err", err))
				return err
			}

			// Declare OTP exchange
			if err := rmq.DeclareExchange(
				cfg.RabbitMQ.ChatIDExchange.Name,
				cfg.RabbitMQ.ChatIDExchange.Kind,
				cfg.RabbitMQ.ChatIDExchange.Durable,
				cfg.RabbitMQ.ChatIDExchange.AutoDeleted,
				cfg.RabbitMQ.ChatIDExchange.Internal,
				cfg.RabbitMQ.ChatIDExchange.NoWait,
				cfg.RabbitMQ.ChatIDExchange.Args.ToMap(),
			); err != nil {
				log.Error("failed to declare OTP exchange", slog.Any("err", err))
			}

			// Declare OTP Queue
			if _, err := rmq.DeclareQueue(
				cfg.RabbitMQ.ChatIDQueue.Name,
				cfg.RabbitMQ.ChatIDQueue.Durable,
				cfg.RabbitMQ.ChatIDQueue.AutoDeleted,
				cfg.RabbitMQ.ChatIDQueue.Exclusive,
				cfg.RabbitMQ.ChatIDQueue.NoWait,
				cfg.RabbitMQ.ChatIDQueue.Args.ToMap(),
			); err != nil {
				log.Error("failed to declare OTP queue", slog.Any("err", err))
			}

			// Bind OTP queue to OTP exchange
			if err := rmq.BindQueueToExchange(
				cfg.RabbitMQ.ChatIDQueue.Name,
				cfg.RabbitMQ.ChatIDExchange.Name,
				cfg.RabbitMQ.ChatIDRoutingKey,
			); err != nil {
				log.Error("failed to bind OTP queue", slog.Any("err", err))
			}

			// start tg bot
			wg.Add(1)
			go func() {
				defer wg.Done()
				telegram.RunTg(
					ctx,
					tgClient,
					cfg.TelegramBot.BatchSize,
					rmq,
					log,
				)
			}()

			// start OTP queue consumer
			consumeOTP := sender.NewConsumeOTP(
				rmq,
				tgClient,
				log,
			)
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := consumeOTP.Start(ctx, cfg.RabbitMQ.OTPQueue.Name); err != nil {
					log.Error("failed to start OTP consumer", slog.Any("err", err))
				}
			}()

			log.Info("tg bot starting at:", slog.Any("port", cfg.Server.Port))
			<-ctx.Done()
			wg.Wait()

			return nil
		}}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}
