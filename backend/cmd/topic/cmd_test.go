package topic

import (
	"bytes"
	"strings"
	"testing"

	"github.com/geoberle/mathilde/backend/model"
	"github.com/geoberle/mathilde/backend/store"
)

func TestSlugify(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := store.TopicSlug(tt.name)
			if got != tt.want {
				t.Errorf("store.TopicSlug(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestParseConceptsJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name: "valid JSON array",
			input: `[
				{"id": "a", "name": "Concept A", "prerequisites": []},
				{"id": "b", "name": "Concept B", "prerequisites": ["a"]}
			]`,
			want: 2,
		},
		{
			name:  "with code fence",
			input: "```json\n" + `[{"id": "a", "name": "A", "prerequisites": []}]` + "\n```",
			want:  1,
		},
		{
			name:    "empty array",
			input:   `[]`,
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			input:   `not json`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			concepts, err := parseConceptsFromText(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(concepts) != tt.want {
				t.Errorf("concept count: got=%d want=%d", len(concepts), tt.want)
			}
			for _, c := range concepts {
				if c.Status != model.StatusNotStarted {
					t.Errorf("concept %q: status got=%q want=%q", c.ID, c.Status, model.StatusNotStarted)
				}
			}
		})
	}
}

func TestBuildProfileContext(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		profile *model.Profile
		want    string
	}{
		{
			name:    "empty profile",
			profile: &model.Profile{},
			want:    "",
		},
		{
			name: "grade only",
			profile: &model.Profile{
				Preferences: map[string]string{"grade": "5. Klasse"},
			},
			want: "- Schulstufe: 5. Klasse",
		},
		{
			name: "mission only",
			profile: &model.Profile{
				Mission: "Bruchrechnung vor der Schularbeit meistern",
			},
			want: "- Lernziel: Bruchrechnung vor der Schularbeit meistern",
		},
		{
			name: "grade and mission",
			profile: &model.Profile{
				Mission:     "Schularbeit bestehen",
				Preferences: map[string]string{"grade": "6. Klasse"},
			},
			want: "- Schulstufe: 6. Klasse\n- Lernziel: Schularbeit bestehen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildProfileContext(tt.profile)
			if got != tt.want {
				t.Errorf("buildProfileContext() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShowOutput(t *testing.T) {
	t.Parallel()
	topic := model.Topic{
		Name: "Bruchrechnung",
		Concepts: []model.Concept{
			{ID: "was-ist-ein-bruch", Name: "Was ist ein Bruch?"},
			{ID: "erweitern", Name: "Erweitern von Brüchen", Prerequisites: []string{"was-ist-ein-bruch"}},
		},
	}

	var buf bytes.Buffer
	printTopicTree(&buf, &topic)
	output := buf.String()

	if !strings.Contains(output, "Bruchrechnung") {
		t.Error("output missing topic name")
	}
	if !strings.Contains(output, "Was ist ein Bruch?") {
		t.Error("output missing concept name")
	}
	if !strings.Contains(output, "Requires: Was ist ein Bruch?") {
		t.Error("output missing prerequisite link")
	}
}

func TestShowOutputUnknownPrerequisite(t *testing.T) {
	t.Parallel()
	topic := model.Topic{
		Name: "Test",
		Concepts: []model.Concept{
			{ID: "a", Name: "Concept A", Prerequisites: []string{"unknown-id"}},
		},
	}

	var buf bytes.Buffer
	printTopicTree(&buf, &topic)
	output := buf.String()

	if !strings.Contains(output, "Requires: unknown-id") {
		t.Errorf("expected raw ID fallback for unknown prerequisite, got:\n%s", output)
	}
}
