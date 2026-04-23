package creds

import "time"

type PortableCreds struct {
	AccessKeyID    string    `json:"access_key_id"`
	SecretAccessKey string    `json:"secret_access_key"`
	SessionToken   string    `json:"session_token"`
	Expiration     time.Time `json:"expiration"`
	Region         string    `json:"region"`
	ProfileName    string    `json:"profile_name"`
}
