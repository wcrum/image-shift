package main

import (
	"testing"
)

func TestImageSwap_SwapImage(t *testing.T) {
	type fields struct {
		Default  string
		Mappings []ImageMapping
		Tests    []string
	}
	type args struct {
		image string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ImageSwap{
				Default:  tt.fields.Default,
				Mappings: tt.fields.Mappings,
				Tests:    tt.fields.Tests,
			}
			if got := i.SwapImage(tt.args.image); got != tt.want {
				t.Errorf("ImageSwap.SwapImage() = %v, want %v", got, tt.want)
			}
		})
	}
}
