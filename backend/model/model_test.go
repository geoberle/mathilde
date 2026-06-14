package model_test

import (
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"context"

	"github.com/geoberle/mathilde/backend/model"
	"github.com/geoberle/mathilde/backend/store"
)

func newTestClient(t *testing.T) *firestore.Client {
	t.Helper()
	if len(os.Getenv("FIRESTORE_EMULATOR_HOST")) == 0 {
		t.Skip("FIRESTORE_EMULATOR_HOST not set — skipping integration test")
	}
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, store.ProjectID())
	if err != nil {
		t.Fatalf("creating firestore client: %v", err)
	}
	return client
}

func TestLearningRecordRoundTrip(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	want := model.LearningRecord{
		ConceptID:     "erweitern",
		Status:        model.StatusMastered,
		Narrative:     "Demonstrated understanding of Erweitern. Initially confused direction...",
		LastPracticed: time.Date(2026, 6, 10, 14, 0, 0, 0, time.UTC),
		NextReview:    time.Date(2026, 6, 17, 14, 0, 0, 0, time.UTC),
		Interval:      7,
		ErrorPatterns: []string{"confuses-with-kuerzen"},
	}

	doc := client.Collection("users").Doc("test-user").Collection("learning-records").Doc("erweitern")
	_, err := doc.Set(ctx, want)
	if err != nil {
		t.Fatalf("writing learning record: %v", err)
	}
	defer doc.Delete(ctx)
	defer client.Close()

	snap, err := doc.Get(ctx)
	if err != nil {
		t.Fatalf("reading learning record: %v", err)
	}

	var got model.LearningRecord
	if err := snap.DataTo(&got); err != nil {
		t.Fatalf("decoding learning record: %v", err)
	}

	if got.ConceptID != want.ConceptID {
		t.Errorf("conceptId: got=%q want=%q", got.ConceptID, want.ConceptID)
	}
	if got.Status != want.Status {
		t.Errorf("status: got=%q want=%q", got.Status, want.Status)
	}
	if got.Narrative != want.Narrative {
		t.Errorf("narrative: got=%q want=%q", got.Narrative, want.Narrative)
	}
	if got.Interval != want.Interval {
		t.Errorf("interval: got=%d want=%d", got.Interval, want.Interval)
	}
	if len(got.ErrorPatterns) != 1 || got.ErrorPatterns[0] != "confuses-with-kuerzen" {
		t.Errorf("errorPatterns: got=%v want=%v", got.ErrorPatterns, want.ErrorPatterns)
	}
}

func TestProfileRoundTrip(t *testing.T) {
	client := newTestClient(t)
	ctx := context.Background()

	want := model.Profile{
		XP:      150,
		Level:   3,
		Mission: "Bruchrechnung vor der Schularbeit meistern",
	}

	doc := client.Collection("users").Doc("test-user").Collection("profile").Doc("main")
	_, err := doc.Set(ctx, want)
	if err != nil {
		t.Fatalf("writing profile: %v", err)
	}
	defer doc.Delete(ctx)
	defer client.Close()

	snap, err := doc.Get(ctx)
	if err != nil {
		t.Fatalf("reading profile: %v", err)
	}

	var got model.Profile
	if err := snap.DataTo(&got); err != nil {
		t.Fatalf("decoding profile: %v", err)
	}

	if got.XP != want.XP {
		t.Errorf("xp: got=%d want=%d", got.XP, want.XP)
	}
	if got.Level != want.Level {
		t.Errorf("level: got=%d want=%d", got.Level, want.Level)
	}
	if got.Mission != want.Mission {
		t.Errorf("mission: got=%q want=%q", got.Mission, want.Mission)
	}
}
