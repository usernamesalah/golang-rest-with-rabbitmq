/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"tenant/internal/container"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process [client-id] [payload]",
	Short: "Process a tenant",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		clientID := args[0]
		payload := args[1]
		sc, err := container.InitSharedComponent()
		if err != nil {
			logrus.Infof("Failed to process tenant: %v\n", err)
			return
		}
		defer sc.DB.Close()
		defer sc.RabbitMQConn.Close()

		cc := container.NewHandlerComponent(sc)

		err = cc.TenantUsecase.ProcessPayload(cmd.Context(), clientID, payload)
		if err != nil {
			logrus.Infof("Failed to process tenant: %v\n", err)
			return
		}

		logrus.Info("Process Tenant successfully")
	},
}

func init() {
	rootCmd.AddCommand(processCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// processCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
