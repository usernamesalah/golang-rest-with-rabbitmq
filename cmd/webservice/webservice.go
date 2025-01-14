package webservice

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tenant/internal/api/http/router"

	"time"

	"tenant/internal/container"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func Run(sc *container.SharedComponent) error {

	cc := container.NewHandlerComponent(sc)

	logrus.Info("Initializing the web server ...")
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
			logrus.WithFields(logrus.Fields{
				"Latency":   v.Latency.String(),
				"Remote IP": c.RealIP(),
				"URI":       v.URI,
				"Method":    c.Request().Method,
				"status":    v.Status,
			}).Info("request")

			return nil
		},
	}))

	e.Validator = &requestValidator{}

	// init route
	router.Init(e, cc, sc)

	// Start server
	server := &http.Server{
		Addr:         "0.0.0.0:" + sc.Conf.Server.Port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	serverErrors := make(chan error, 1)
	// mulai listening server
	go func() {
		logrus.Infof("server listening on %v", server.Addr)
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
		logrus.Info("caught signal, shutting down")

		// Jika ada shutdown, meminta tambahan waktu 10 detik untuk menyelesaikan proses yang sedang berjalan.
		const timeout = 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		defer sc.DB.Close()
		if err := sc.RabbitMQConn.Close(); err != nil {
			logrus.Errorf("error: gracefully shutting down rabbitmq connection : %s", err)
		}

		if err := server.Shutdown(ctx); err != nil {
			logrus.Errorf("error: gracefully shutting down server: %s", err)
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
