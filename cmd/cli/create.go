/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"tenant/internal/container"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new tenant",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		sc, err := container.InitSharedComponent()
		if err != nil {
			logrus.Infof("Failed to create tenant: %v\n", err)
			return
		}
		defer sc.DB.Close()
		defer sc.RabbitMQConn.Close()

		cc := container.NewHandlerComponent(sc)

		tenant, err := cc.TenantUsecase.CreateTenant(cmd.Context(), name)
		if err != nil {
			logrus.Infof("Failed to create tenant: %v\n", err)
			return
		}

		logrus.Infof("Tenant created successfully: %+v\n", tenant.ClientID)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
