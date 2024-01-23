package main

import (
	"alfred/backend/initial"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "server",
		Short: "accounts server.",
		Long:  `accounts server.`,
		Run:   initial.StartServer,
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "migrate-db",
		Short: "Auto migrate database by gorm.",
		Long:  `Auto migrate database by gorm.`,
		Run:   initial.MigrateDB,
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("exec cmd error:", err.Error())
		os.Exit(1)
	}
}
