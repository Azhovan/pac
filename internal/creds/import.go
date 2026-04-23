package creds

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
)

func Import(filePath string, profileOverride string) (*PortableCreds, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file %s: %w", filePath, err)
	}

	var pc PortableCreds
	if err := json.Unmarshal(data, &pc); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	profileName := pc.ProfileName
	if profileOverride != "" {
		profileName = profileOverride
	}

	if pc.Expiration.Before(time.Now()) {
		fmt.Fprintf(os.Stderr, "Warning: credentials expired at %s\n", pc.Expiration.Local().Format("2006-01-02 15:04:05"))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to determine home directory: %w", err)
	}

	awsDir := filepath.Join(home, ".aws")
	if err := os.MkdirAll(awsDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create %s: %w", awsDir, err)
	}

	if err := writeCredentials(awsDir, profileName, &pc); err != nil {
		return nil, err
	}

	if pc.Region != "" {
		if err := writeConfig(awsDir, profileName, &pc); err != nil {
			return nil, err
		}
	}

	return &pc, nil
}

func writeCredentials(awsDir string, profileName string, pc *PortableCreds) error {
	credPath := filepath.Join(awsDir, "credentials")

	credFile, err := ini.LooseLoad(credPath)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", credPath, err)
	}

	section := credFile.Section(profileName)
	section.Key("aws_access_key_id").SetValue(pc.AccessKeyID)
	section.Key("aws_secret_access_key").SetValue(pc.SecretAccessKey)
	section.Key("aws_session_token").SetValue(pc.SessionToken)

	if err := credFile.SaveTo(credPath); err != nil {
		return fmt.Errorf("failed to write %s: %w", credPath, err)
	}
	if err := os.Chmod(credPath, 0600); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", credPath, err)
	}

	return nil
}

func writeConfig(awsDir string, profileName string, pc *PortableCreds) error {
	configPath := filepath.Join(awsDir, "config")

	cfgFile, err := ini.LooseLoad(configPath)
	if err != nil {
		return fmt.Errorf("failed to load %s: %w", configPath, err)
	}

	sectionName := fmt.Sprintf("profile %s", profileName)
	if profileName == "default" {
		sectionName = "default"
	}

	section := cfgFile.Section(sectionName)
	section.Key("region").SetValue(pc.Region)

	if err := cfgFile.SaveTo(configPath); err != nil {
		return fmt.Errorf("failed to write %s: %w", configPath, err)
	}
	if err := os.Chmod(configPath, 0600); err != nil {
		return fmt.Errorf("failed to set permissions on %s: %w", configPath, err)
	}

	return nil
}
