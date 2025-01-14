package container

import (
	"tenant/infrastructure/config"
	"tenant/infrastructure/database"
	"tenant/pkg/logger"
	"time"

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

func InitSharedComponent() (*SharedComponent, error) {

	// Load configuration
	config.Init()
	conf := config.Get()

	// Start Logger
	log := logger.NewLogger(*conf)

	// Start Database
	database, err := database.InitializeDatabase(conf)
	if err != nil {
		log.Errorf("web failed to init db %v", err)
		return nil, err
	}

	mqConn, err := amqp091.Dial(conf.RabbitMQ.URL)
	if err != nil {
		log.Errorf("failed to connect to RabbitMQ: %v", err)
		return nil, err
	}

	go func() {
		for {
			<-mqConn.NotifyClose(make(chan *amqp091.Error))
			for i := 0; i < conf.RabbitMQ.MaxReconnects; i++ {
				if mqConnr, err := amqp091.Dial(conf.RabbitMQ.URL); err == nil {
					log.Println("Reconnected to RabbitMQ")
					mqConn = mqConnr
					break
				}
				log.Printf("Reconnection attempt %d failed, retrying in %dms...\n", i+1, conf.RabbitMQ.ReconnectDelay)
				time.Sleep(time.Duration(conf.RabbitMQ.ReconnectDelay) * time.Millisecond)
			}
		}
	}()

	sharedComponent := &SharedComponent{
		DB:           database,
		Conf:         conf,
		Log:          log,
		RabbitMQConn: mqConn,
	}

	return sharedComponent, nil
}
