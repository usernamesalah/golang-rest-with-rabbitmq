package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tenant/infrastructure/config"
	"tenant/infrastructure/database"
	"tenant/internal/api/http/router"
	"tenant/pkg/logger"

	"time"

	"tenant/internal/container"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

// @title Api Documentation for tenant apps backend
// @version 0.1
// @description API documentation for tenant apps backend

// @contact.name Tenant Apps
// @contact.email no-reply@b2b-tenant.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
func main() {

	// Load configuration
	config.Init()
	conf := config.Get()
	log := logger.NewLogger(*conf)

	if err := run(log); err != nil {
		log.Errorf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run(log *logrus.Logger) error {
	conf := config.Get()

	// Start Database
	database, err := database.InitializeDatabase(conf)
	if err != nil {
		log.Error("web failed to init db", zap.Error(err))
		return err
	}

	mqConn, err := amqp091.Dial(conf.RabbitMQ.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	go func() {
		for {
			<-mqConn.NotifyClose(make(chan *amqp091.Error))
			log.Println("Connection lost, attempting to reconnect...")
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

	sharedComponent := &container.SharedComponent{
		DB:           database,
		Conf:         conf,
		Log:          log,
		RabbitMQConn: mqConn,
	}

	cc := container.NewHandlerComponent(sharedComponent)

	log.Info("Initializing the web server ...")
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Header.Set("Cache-Control", "max-age:3600, public")
			return next(c)
		}
	})

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Info("request",
				zap.String("Latency", v.Latency.String()),
				zap.String("Remote IP", c.RealIP()),
				zap.String("URI", v.URI),
				zap.String("Method", c.Request().Method),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	e.Validator = &requestValidator{}

	// init route
	router.Init(e, cc, sharedComponent)

	// Start server
	server := &http.Server{
		Addr:         "0.0.0.0:" + conf.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	serverErrors := make(chan error, 1)
	// mulai listening server
	go func() {
		log.Info("server listening on", zap.String("address", server.Addr))
		serverErrors <- e.StartServer(server)
	}()

	// Membuat channel untuk mendengarkan sinyal interupsi/terminate dari OS.
	// Menggunakan channel buffered karena paket signal membutuhkannya.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Mengontrol penerimaan data dari channel,
	// jika ada error saat listenAndServe server maupun ada sinyal shutdown yang diterima
	select {
	case err := <-serverErrors:
		return fmt.Errorf("starting server: %v", err)

	case <-shutdown:
		log.Info("caught signal, shutting down")

		// Jika ada shutdown, meminta tambahan waktu 10 detik untuk menyelesaikan proses yang sedang berjalan.
		const timeout = 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		database.Close()
		if err := mqConn.Close(); err != nil {
			log.Errorf("error: gracefully shutting down rabbitmq connection : %s", err)
		}

		if err := server.Shutdown(ctx); err != nil {
			log.Errorf("error: gracefully shutting down server: %s", err)
			if err := server.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %v", err)
			}
		}

	}

	return nil
}

type requestValidator struct{}

func (rv *requestValidator) Validate(i interface{}) (err error) {
	_, err = govalidator.ValidateStruct(i)
	return
}
