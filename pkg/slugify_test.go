package pkg_test

import (
	"testing"

	"github.com/0xdod/go-realworld/pkg"
)

func TestSlugify(t *testing.T) {
	cases := []struct {
		Name     string
		Input    string
		Expected string
	}{
		{
			"slug correctly",
			"How to train your dragon",
			"how-to-train-your-dragon",
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got := pkg.Slugify(tc.Input)

			if got != tc.Expected {
				t.Errorf("expected %q, but got %q", tc.Expected, got)
			}
		})
	}
}
