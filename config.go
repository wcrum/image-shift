package main

import (
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"
)

type config struct {
	certFile string
	keyFile  string

	// probably need to extend this
}

func initEnv() *config {
	cfg := &config{}

	if certFile := os.Getenv("TLS_CERT_FILE"); certFile != "" {
		cfg.certFile = certFile
	}

	if keyFile := os.Getenv("TLS_KEY_FILE"); keyFile != "" {
		cfg.keyFile = keyFile
	}

	return cfg
}

func initConfig() *ImageSwapConfig {
	// default location will be /etc/imageswap-config/imageswap.yaml
	imageSwapFilePath := "/etc/imageswap-configmap/imageswap.yaml"
	cfg := &ImageSwapConfig{}

	content, err := os.ReadFile(imageSwapFilePath)
	if err != nil {
	}

	err = yaml.Unmarshal(content, cfg)

	return cfg
}
