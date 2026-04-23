package creds

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
)

func Export(ctx context.Context, profile string, outputPath string) (*PortableCreds, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithSharedConfigProfile(profile),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load profile %q: %w (is it defined in ~/.aws/config?)", profile, err)
	}

	credentials, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credentials for profile %q: %w (have you run 'aws sso login --profile %s'?)", profile, err, profile)
	}

	pc := &PortableCreds{
		AccessKeyID:    credentials.AccessKeyID,
		SecretAccessKey: credentials.SecretAccessKey,
		SessionToken:   credentials.SessionToken,
		Expiration:     credentials.Expires,
		Region:         cfg.Region,
		ProfileName:    profile,
	}

	data, err := json.MarshalIndent(pc, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal credentials: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return nil, fmt.Errorf("failed to write credentials to %s: %w", outputPath, err)
	}

	return pc, nil
}
