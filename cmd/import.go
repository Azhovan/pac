package cmd

import (
	"fmt"

	"github.com/jabar/aws-token-persister/internal/creds"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import credentials from a portable JSON file into ~/.aws/credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		profile, _ := cmd.Flags().GetString("profile")

		pc, err := creds.Import(file, profile)
		if err != nil {
			return err
		}

		targetProfile := pc.ProfileName
		if profile != "" {
			targetProfile = profile
		}

		fmt.Printf("Credentials imported successfully\n")
		fmt.Printf("  Profile:    %s\n", targetProfile)
		fmt.Printf("  Region:     %s\n", pc.Region)
		fmt.Printf("  Expires at: %s\n", pc.Expiration.Local().Format("2006-01-02 15:04:05"))
		return nil
	},
}

func init() {
	importCmd.Flags().StringP("file", "f", "", "Path to the portable credentials JSON file (required)")
	importCmd.Flags().StringP("profile", "p", "", "Override the profile name to write (defaults to profile from the file)")
	_ = importCmd.MarkFlagRequired("file")
}
