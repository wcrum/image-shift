package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Define the structure of the imageswap configuration
type ImageSwapConfig struct {
	ImageSwap ImageSwap
}
type ImageSwap struct {
	Default  string         `yaml:"default"`
	Mappings []ImageMapping `yaml:"mappings"`
	Tests    []string       `yaml:"tests,flow"`
}

// Define the structure of an image mapping
type ImageMapping struct {
	Registry string `yaml:"registry"`
	Action   string `yaml:"action"`
	Target   string `yaml:"target"`
}

func (i *ImageSwap) SwapImage(image string) string {
	ref, _ := name.ParseReference(image, name.WithDefaultRegistry(i.Default))

	registry := ref.Context().RegistryStr()

	// if registry == default registry return image
	if string(registry) == i.Default {
		return ref.Name()
	}

	var newImage string

	for _, m := range i.Mappings {
		switch m.Action {
		case "ignore":
			if m.Registry == registry {
				log.Info(fmt.Sprintf("Ignoring %s.", ref.Name()))
				return ref.Name()
			}
		case "swap":
			if m.Registry == registry {
				identifier := ref.Identifier()
				switch len(strings.Split(identifier, ":")) {
				case 1:
					newImage = fmt.Sprintf("%s/%s:%s", m.Target, ref.Context().RepositoryStr(), ref.Identifier())
				case 2:
					newImage = fmt.Sprintf("%s/%s@%s", m.Target, ref.Context().RepositoryStr(), ref.Identifier())
				default:
					newImage = fmt.Sprintf("%s/%s", m.Target, ref.Context().RepositoryStr())
				}
			}
		}
	}

	return newImage
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	var config ImageSwapConfig

	path := "imageswap.yaml"

	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(content, &config)

	for _, image := range config.ImageSwap.Tests {
		fmt.Printf("%s -> %s\n", image, config.ImageSwap.SwapImage(image))
	}
}
