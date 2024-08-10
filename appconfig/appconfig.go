package appconfig

import (
	"os"
	"strings"
)

type AppConfig struct {
	IsDevelopment bool
	Remote        string
	Host          string
	WsUrl         string
	UploadUrl     string
}

func (config *AppConfig) init() {
	config.IsDevelopment = false

	if _, err := os.Stat(".env"); err == nil {
		content, err := os.ReadFile(".env")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(strings.TrimPrefix(line, "\uFEFF")) // Remove BOM if present
				if strings.HasPrefix(line, "ENV=") {
					envValue := strings.TrimPrefix(line, "ENV=")
					config.IsDevelopment = strings.EqualFold(strings.TrimSpace(envValue), "development")
					break
				}
			}
		}
	}

	if config.IsDevelopment {
		config.Remote = "localhost:8899"
		config.Host = "http://" + config.Remote
		config.WsUrl = "ws://" + config.Remote + "/x/%s/ws"
	} else {
		config.Remote = "1fps.video"
		config.Host = "https://" + config.Remote
		config.WsUrl = "wss://" + config.Remote + "/x/%s/ws"
	}

	config.UploadUrl = config.Host + "/upload"
}

func New() *AppConfig {
	config := &AppConfig{}
	config.init()
	return config
}
