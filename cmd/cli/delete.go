/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cli

import (
	"tenant/internal/container"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [client-id]",
	Short: "Delete a tenant",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		clientID := args[0]
		sc, err := container.InitSharedComponent()
		if err != nil {
			logrus.Infof("Failed to delete tenant: %v\n", err)
			return
		}
		defer sc.DB.Close()
		defer sc.RabbitMQConn.Close()

		cc := container.NewHandlerComponent(sc)

		err = cc.TenantUsecase.DeleteTenant(cmd.Context(), clientID)
		if err != nil {
			logrus.Infof("Failed to delete tenant: %v\n", err)
			return
		}

		logrus.Info("Tenant deleted successfully")
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
