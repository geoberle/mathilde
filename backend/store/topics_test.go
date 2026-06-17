package store

import "testing"

func Test_slugify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want string
	}{
		{name: "Bruchrechnung", want: "bruchrechnung"},
		{name: "Gleichungen lösen", want: "gleichungen-loesen"},
		{name: "Flächen und Körper", want: "flaechen-und-koerper"},
		{name: "Größen und Einheiten", want: "groessen-und-einheiten"},
		{name: "Maße und Maßstäbe", want: "masse-und-massstaebe"},
		{name: "Flächen & Körper", want: "flaechen-koerper"},
		{name: "  spaces  ", want: "spaces"},
		{name: "a--b", want: "a-b"},
		{name: "ABC", want: "abc"},
		{name: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := slugify(tt.name)
			if got != tt.want {
				t.Errorf("slugify(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
