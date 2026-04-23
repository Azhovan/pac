package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "pac",
	Short: "PAC - Portable AWS Credentials. Export and import AWS SSO credentials between machines.",
}

func init() {
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
