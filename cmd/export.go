package cmd

import (
	"fmt"

	"github.com/Azhovan/pac/internal/creds"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export AWS SSO credentials to a portable JSON file",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		output, _ := cmd.Flags().GetString("output")

		pc, err := creds.Export(cmd.Context(), profile, output)
		if err != nil {
			return err
		}

		fmt.Printf("Credentials exported to %s\n", output)
		fmt.Printf("  Profile:    %s\n", pc.ProfileName)
		fmt.Printf("  Region:     %s\n", pc.Region)
		fmt.Printf("  Expires at: %s\n", pc.Expiration.Local().Format("2006-01-02 15:04:05"))
		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("profile", "p", "", "AWS SSO profile name (required)")
	exportCmd.Flags().StringP("output", "o", "aws-creds.json", "Output file path")
	_ = exportCmd.MarkFlagRequired("profile")
}
