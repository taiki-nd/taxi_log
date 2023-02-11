package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	LogFile           string
	Sql               string
	Host              string
	Port              string
	Name              string
	User              string
	Password          string
	Url               string
	GcsBucketName     string
	GcsObjectPath     string
	GcsKeyPath        string
	StripePublicOkKey string
	StripeSecretKey   string
	DiscordWebhookUrl string
}

var Config ConfigList

func init() {
	// iniファイルの読み込み
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Printf("failed to load config.ini: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		LogFile:           cfg.Section("taxi_log").Key("log_file").String(),
		Sql:               cfg.Section("db").Key("sql").String(),
		Host:              cfg.Section("db").Key("host").String(),
		Port:              cfg.Section("db").Key("port").String(),
		Name:              cfg.Section("db").Key("name").String(),
		User:              cfg.Section("db").Key("user").String(),
		Password:          cfg.Section("db").Key("password").String(),
		Url:               cfg.Section("cors").Key("url").String(),
		GcsBucketName:     cfg.Section("gcp").Key("gcs_bucket_name").String(),
		GcsObjectPath:     cfg.Section("gcp").Key("gcs_object_path").String(),
		GcsKeyPath:        cfg.Section("gcp").Key("gcs_key_path").String(),
		StripePublicOkKey: cfg.Section("stripe").Key("public_ok_key").String(),
		StripeSecretKey:   cfg.Section("stripe").Key("secret_key").String(),
		DiscordWebhookUrl: cfg.Section("discord").Key("webhook_url").String(),
	}
}
