package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
)

// Define the structure of the imageswap configuration
type ImageSwapConfig struct {
	ImageSwap ImageSwap `yaml:"imageswap"`
}
type ImageSwap struct {
	Default  string       `yaml:"default"`
	Mappings ImageMapping `yaml:"mappings"`
	Tests    []string     `yaml:"tests,flow"`
}

// Define the structure of an image mapping
type ImageMapping struct {
	Swap       []Swap      `yaml:"swap"`
	ExactSwap  []ExactSwap `yaml:"exact-swap"`
	RegexSwap  []RegexSwap `yaml:"regex-swap"`
	LegacySwap []string    `yaml:"legacy-swap"`
}

// Swap represents the "swap" mapping in the YAML.
type Swap struct {
	Registry string `yaml:"registry"`
	Target   string `yaml:"target"`
}

// ExactSwap represents the "exact-swap" mapping in the YAML.
type ExactSwap struct {
	Image  string `yaml:"image"`
	Target string `yaml:"target"`
}

// RegexSwap represents the "regex-swap" mapping in the YAML.
type RegexSwap struct {
	Expression string `yaml:"expression"`
	Target     string `yaml:"target"`
}

func (i *ImageSwap) SwapImage(image string) string {
	ref, _ := name.ParseReference(image, name.WithDefaultRegistry(i.Default))

	registry := ref.Context().RegistryStr()

	// if registry == default registry return image
	var newImage string

	if string(registry) == i.Default {
		newImage = ref.Name()
	}

	for _, swap := range i.Mappings.Swap {
		if swap.Registry == registry {
			identifier := ref.Identifier()
			switch len(strings.Split(identifier, ":")) {
			case 1:
				newImage = fmt.Sprintf("%s/%s:%s", swap.Target, ref.Context().RepositoryStr(), ref.Identifier())
			case 2:
				tag := strings.Split(strings.Split(image, "@")[0], ":")[1]
				newImage = fmt.Sprintf("%s/%s:%s@%s", swap.Target, ref.Context().RepositoryStr(), tag, ref.Identifier())
			default:
				newImage = fmt.Sprintf("%s/%s", swap.Target, ref.Context().RepositoryStr())
			}
			break
		}
	}

	for _, exactSwap := range i.Mappings.ExactSwap {
		if exactSwap.Image == ref.String() {
			newImage = exactSwap.Target
			break
		}
	}

	for _, regexSwap := range i.Mappings.RegexSwap {
		re := regexp.MustCompile(regexSwap.Expression)

		match := re.FindStringSubmatch(ref.String())
		if match != nil {
			newRepository := strings.Replace(ref.String(), match[0], "", -1)

			newImage = regexSwap.Target

			if len(match) > 1 {
				for m := 1; m < len(match); m++ {
					newImage = strings.Replace(newImage, "$"+strconv.Itoa(m), match[m], -1)
					newImage = fmt.Sprintf("%s%s", newImage, newRepository)
				}
			}
			break
		}
	}

	return newImage
}
