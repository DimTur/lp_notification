package notification

import (
	"context"
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
			rmq, err := initRabbitMQ(cfg)
			if err != nil {
				log.Error("failed init rabbit mq", slog.Any("err", err))
			}

			// Declare OTP exchange
			if err := declareOTPExchange(rmq, cfg); err != nil {
				log.Error("failed to declare OTP exchange", slog.Any("err", err))
			}

			// Declare ChatID exchange
			if err := declareChatIDExchange(rmq, cfg); err != nil {
				log.Error("failed to declare ChatID exchange", slog.Any("err", err))
			}

			// Declare and bind OTP Queue
			if err := declareQueueAndBind(
				rmq,
				cfg.RabbitMQ.OTP.OTPQueue,
				cfg.RabbitMQ.OTP.OTPExchange.Name,
				cfg.RabbitMQ.OTP.OTPRoutingKey,
			); err != nil {
				log.Error("failed to declare and bind OTP Queue", slog.Any("err", err))
			}

			// Declare and bind ChatID Queue
			if err := declareQueueAndBind(
				rmq,
				cfg.RabbitMQ.Chat.ChatIDQueue,
				cfg.RabbitMQ.Chat.ChatIDExchange.Name,
				cfg.RabbitMQ.Chat.ChatIDRoutingKey,
			); err != nil {
				log.Error("failed to declare and bind ChatID Queue", slog.Any("err", err))
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

			startConsumers(ctx, cfg, rmq, tgClient, log, &wg)

			log.Info("tg bot starting at:", slog.Any("port", cfg.Server.Port))
			<-ctx.Done()
			wg.Wait()

			return nil
		}}

	c.Flags().StringVar(&configPath, "config", "", "path to config")
	return c
}

func initRabbitMQ(cfg *config.Config) (*rabbitmq_store.RMQClient, error) {
	rmqUrl := fmt.Sprintf(
		"amqp://%s:%s@%s:%d/",
		cfg.RabbitMQ.UserName,
		cfg.RabbitMQ.Password,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port,
	)
	return rabbitmq_store.NewClient(rmqUrl)
}

func startConsumers(
	ctx context.Context,
	cfg *config.Config,
	rmq *rabbitmq_store.RMQClient,
	tgClient *tgclient.TgClient,
	log *slog.Logger,
	wg *sync.WaitGroup,
) {
	otpConsumer := sender.NewConsumeOTP(rmq, tgClient, log)
	shareConsumer := sender.NewConsumeNotification(rmq, tgClient, log)

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := otpConsumer.Start(
			ctx,
			cfg.RabbitMQ.OTP.OTPConsumer.Queue,
			cfg.RabbitMQ.OTP.OTPConsumer.Consumer,
			cfg.RabbitMQ.OTP.OTPConsumer.AutoAck,
			cfg.RabbitMQ.OTP.OTPConsumer.Exclusive,
			cfg.RabbitMQ.OTP.OTPConsumer.NoLocal,
			cfg.RabbitMQ.OTP.OTPConsumer.NoWait,
			cfg.RabbitMQ.OTP.OTPConsumer.ConsumerArgs.ToMap(),
		); err != nil {
			log.Error("failed to start otp consumer", slog.Any("err", err))
		}
	}()

	go func() {
		defer wg.Done()
		if err := shareConsumer.Start(
			ctx,
			cfg.RabbitMQ.Notification.NotificationConsumer.Queue,
			cfg.RabbitMQ.Notification.NotificationConsumer.Consumer,
			cfg.RabbitMQ.Notification.NotificationConsumer.AutoAck,
			cfg.RabbitMQ.Notification.NotificationConsumer.Exclusive,
			cfg.RabbitMQ.Notification.NotificationConsumer.NoLocal,
			cfg.RabbitMQ.Notification.NotificationConsumer.NoWait,
			cfg.RabbitMQ.Notification.NotificationConsumer.ConsumerArgs.ToMap(),
		); err != nil {
			log.Error("failed to start notification consumer", slog.Any("err", err))
		}
	}()
}

func declareQueueAndBind(rmq *rabbitmq_store.RMQClient, queueConfig config.QueueConfig, exchangeName, routingKey string) error {
	// Announcement of the queue
	if _, err := rmq.DeclareQueue(
		queueConfig.Name,
		queueConfig.Durable,
		queueConfig.AutoDeleted,
		queueConfig.Exclusive,
		queueConfig.NoWait,
		queueConfig.Args.ToMap(),
	); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueConfig.Name, err)
	}

	// Binding a queue to an exchange
	if err := rmq.BindQueueToExchange(
		queueConfig.Name,
		exchangeName,
		routingKey,
	); err != nil {
		return fmt.Errorf("failed to bind queue %s to exchange %s: %w", queueConfig.Name, exchangeName, err)
	}

	return nil
}

func declareOTPExchange(rmq *rabbitmq_store.RMQClient, cfg *config.Config) error {
	return rmq.DeclareExchange(
		cfg.RabbitMQ.OTP.OTPExchange.Name,
		cfg.RabbitMQ.OTP.OTPExchange.Kind,
		cfg.RabbitMQ.OTP.OTPExchange.Durable,
		cfg.RabbitMQ.OTP.OTPExchange.AutoDeleted,
		cfg.RabbitMQ.OTP.OTPExchange.Internal,
		cfg.RabbitMQ.OTP.OTPExchange.NoWait,
		cfg.RabbitMQ.OTP.OTPExchange.Args.ToMap(),
	)
}

func declareChatIDExchange(rmq *rabbitmq_store.RMQClient, cfg *config.Config) error {
	return rmq.DeclareExchange(
		cfg.RabbitMQ.Chat.ChatIDExchange.Name,
		cfg.RabbitMQ.Chat.ChatIDExchange.Kind,
		cfg.RabbitMQ.Chat.ChatIDExchange.Durable,
		cfg.RabbitMQ.Chat.ChatIDExchange.AutoDeleted,
		cfg.RabbitMQ.Chat.ChatIDExchange.Internal,
		cfg.RabbitMQ.Chat.ChatIDExchange.NoWait,
		cfg.RabbitMQ.Chat.ChatIDExchange.Args.ToMap(),
	)
}
