package config

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultConfigPath = ".env"

func Load(path string) (*AppConfig, error) {
	if path == "" {
		path = defaultConfigPath
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %w", err)
	}
	defer file.Close()

	env := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	port, _ := strconv.Atoi(env["SERVER_PORT"])
	timeout, _ := strconv.Atoi(env["CLIENT_TIMEOUT_SECONDS"])
	enableRequestID, _ := strconv.ParseBool(env["LOG_ENABLE_REQUEST_ID"])

	cfg := &AppConfig{
		Environment: env["ENVIRONMENT"],
		Server: ServerConfig{
			Host: env["SERVER_HOST"],
			Port: port,
		},
		Client: ClientConfig{
			ServerAddress:  env["CLIENT_SERVER_ADDRESS"],
			TimeoutSeconds: timeout,
		},
		Log: LogConfig{
			EnableRequestID: enableRequestID,
		},
	}

	return cfg, nil
}
