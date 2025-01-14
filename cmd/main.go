package main

import (
	"os"

	"tenant/cmd/cli"
	"tenant/cmd/webservice"
	"tenant/internal/container"

	"github.com/sirupsen/logrus"
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

	sc, err := container.InitSharedComponent()
	if err != nil {
		logrus.Infof("Failed to create tenant: %v\n", err)
		return
	}
	defer sc.DB.Close()
	defer sc.RabbitMQConn.Close()

	// Jalankan CLI atau REST API
	if len(os.Args) > 1 {
		cli.Execute()
	} else {
		if err := webservice.Run(sc); err != nil {
			logrus.Errorf("Error: %v", err)
			os.Exit(1)
		}
	}
}
