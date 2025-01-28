package main

import (
	"testing"
)

func TestImageSwap_SwapImage(t *testing.T) {
	type fields struct {
		Default  string
		Mappings ImageMapping
	}
	tests := []struct {
		name   string
		fields fields
		images []struct {
			input string
			want  string
		}
	}{
		{
			name: "single mapping with multiple images",
			fields: fields{
				Mappings: ImageMapping{
					Swap: []Swap{
						{
							Registry: "gcr.io",
							Target:   "registry.testing.com",
						},
					},
				},
			},
			images: []struct {
				input string
				want  string
			}{
				{input: "gcr.io/testing/test:latest", want: "registry.testing.com/testing/test:latest"},
				{input: "gcr.io/another/image:tag", want: "registry.testing.com/another/image:tag"},
				{input: "gcr.io/example/app:v1", want: "registry.testing.com/example/app:v1"},
				{input: "gcr.io/example/app:tag@sha256:98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624", want: "registry.testing.com/example/app:tag@sha256:98706f0f213dbd440021993a82d2f70451a73698315370ae8615cc468ac06624"},
				{input: "example/app:v1", want: "example/app:v1"},
			},
		},
		{
			name: "single mapping with multiple images and default",
			fields: fields{
				Default: "registry.testing.com",
				Mappings: ImageMapping{
					Swap: []Swap{
						{
							Registry: "gcr.io",
							Target:   "registry.testing.com",
						},
					},
				},
			},
			images: []struct {
				input string
				want  string
			}{
				{input: "example/app:v1", want: "registry.testing.com/example/app:v1"},
			},
		},
		{
			name: "regular expression matching",
			fields: fields{
				Mappings: ImageMapping{
					RegexSwap: []RegexSwap{
						{
							Expression: "gcr.io/([a-zA-Z]*)",
							Target:     "registry.testing.com/cache/$1",
						},
					},
				},
			},
			images: []struct {
				input string
				want  string
			}{
				{input: "gcr.io/company-pepsi/image:latest", want: "registry.testing.com/cache/company-pepsi/image:latest"},
			},
		},
		{
			name: "exact expression matching",
			fields: fields{
				Mappings: ImageMapping{
					ExactSwap: []ExactSwap{
						{
							Image:  "gcr.io/testing/testing:latest",
							Target: "registry.testing.com/abc/abc:test",
						},
					},
				},
			},
			images: []struct {
				input string
				want  string
			}{
				{input: "gcr.io/testing/testing:latest", want: "registry.testing.com/abc/abc:test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ImageSwap{
				Default:  tt.fields.Default,
				Mappings: tt.fields.Mappings,
			}
			for _, img := range tt.images {
				if got := i.SwapImage(img.input); got != img.want {
					t.Errorf("ImageSwap.SwapImage() = %v, want %v", got, img.want)
				}
			}
		})
	}
}
