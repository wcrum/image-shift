package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	log "github.com/sirupsen/logrus"
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
	Image    string `yaml:"image"`
}

func (i *ImageSwap) SwapImage(image string) string {
	ref, _ := name.ParseReference(image, name.WithDefaultRegistry(i.Default))

	registry := ref.Context().RegistryStr()

	// if registry == default registry return image
	var newImage string

	if string(registry) == i.Default {
		newImage = ref.Name()
	}

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
		case "exact-swap":
			if m.Image == ref.String() {
				newImage = m.Target
			}
		case "regex-swap":
			fmt.Println(m.Registry)

			re := regexp.MustCompile(m.Registry)

			match := re.FindStringSubmatch(ref.String())
			if match != nil {
				newRepository := strings.Replace(ref.String(), match[0], "", -1)

				newImage = m.Target

				if len(match) > 1 {
					for m := 1; m < len(match); m++ {
						newImage = strings.Replace(newImage, "$"+strconv.Itoa(m), match[m], -1)
						newImage = fmt.Sprintf("%s%s", newImage, newRepository)
					}
				}
			}

		}
	}

	return newImage
}
