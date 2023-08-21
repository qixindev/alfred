package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "server",
		Short: "accounts server.",
		Long:  `accounts server.`,
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}
	migrateDbCmd = &cobra.Command{
		Use:   "migrate-db",
		Short: "Auto migrate database by gorm.",
		Long:  `Auto migrate database by gorm.`,
		Run: func(cmd *cobra.Command, args []string) {
			migrateDB()
		},
	}
	initFirstCmd = &cobra.Command{
		Use:   "init",
		Short: "Init tenant.",
		Long:  `Init tenant.`,
		Run: func(cmd *cobra.Command, args []string) {
			initFirstRun()
		},
	}
)

func init() {
	rootCmd.AddCommand(migrateDbCmd)
	rootCmd.AddCommand(initFirstCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
