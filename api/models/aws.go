package models

type APIKEYS struct {
	AccessKey       string `json:"access_key"`
	SecretAccessKey string `json:"secret_key"`
	Region          string `json:"region"`
}
