package container

import (
	"tenant/infrastructure/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

// SharedComponent holds shared dependencies between components
type SharedComponent struct {
	Conf         *config.Config
	Log          *logrus.Logger
	DB           *pgxpool.Pool
	RabbitMQConn *amqp091.Connection
}
