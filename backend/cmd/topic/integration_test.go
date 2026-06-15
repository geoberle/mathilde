package topic

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

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

func deleteAllTopics(t *testing.T, client *firestore.Client, uid string) {
	t.Helper()
	ctx := context.Background()
	iter := client.Collection("users").Doc(uid).Collection("topics").Documents(ctx)
	for {
		snap, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			t.Fatalf("listing topics for cleanup: %v", err)
		}
		if _, err := snap.Ref.Delete(ctx); err != nil {
			t.Fatalf("deleting topic %s: %v", snap.Ref.ID, err)
		}
	}
}

func TestTopicRoundTrip(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	uid := "test-topic-user"

	deleteAllTopics(t, client, uid)
	defer deleteAllTopics(t, client, uid)

	want := model.Topic{
		Name: "Bruchrechnung",
		Concepts: []model.Concept{
			{ID: "was-ist-ein-bruch", Name: "Was ist ein Bruch?", Prerequisites: nil},
			{ID: "erweitern", Name: "Erweitern von Brüchen", Prerequisites: []string{"was-ist-ein-bruch"}},
			{ID: "kuerzen", Name: "Kürzen von Brüchen", Prerequisites: []string{"was-ist-ein-bruch"}},
			{ID: "gleichnamig-machen", Name: "Gleichnamig machen", Prerequisites: []string{"erweitern", "kuerzen"}},
		},
		AddedAt: time.Date(2026, 6, 14, 10, 0, 0, 0, time.UTC),
	}

	docID := store.TopicSlug(want.Name)
	doc := client.Collection("users").Doc(uid).Collection("topics").Doc(docID)
	if _, err := doc.Set(ctx, want); err != nil {
		t.Fatalf("writing topic: %v", err)
	}

	snap, err := doc.Get(ctx)
	if err != nil {
		t.Fatalf("reading topic: %v", err)
	}

	var got model.Topic
	if err := snap.DataTo(&got); err != nil {
		t.Fatalf("decoding topic: %v", err)
	}

	if got.Name != want.Name {
		t.Errorf("name: got=%q want=%q", got.Name, want.Name)
	}
	if len(got.Concepts) != len(want.Concepts) {
		t.Fatalf("concepts count: got=%d want=%d", len(got.Concepts), len(want.Concepts))
	}
	for i, wc := range want.Concepts {
		gc := got.Concepts[i]
		if gc.ID != wc.ID {
			t.Errorf("concept[%d].id: got=%q want=%q", i, gc.ID, wc.ID)
		}
		if gc.Name != wc.Name {
			t.Errorf("concept[%d].name: got=%q want=%q", i, gc.Name, wc.Name)
		}
		if len(gc.Prerequisites) != len(wc.Prerequisites) {
			t.Errorf("concept[%d].prerequisites: got=%v want=%v", i, gc.Prerequisites, wc.Prerequisites)
		}
	}
}

func TestTopicOverwrite(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	uid := "test-topic-user"

	deleteAllTopics(t, client, uid)
	defer deleteAllTopics(t, client, uid)

	doc := client.Collection("users").Doc(uid).Collection("topics").Doc("bruchrechnung")

	first := model.Topic{
		Name: "Bruchrechnung",
		Concepts: []model.Concept{
			{ID: "a", Name: "Concept A"},
		},
		AddedAt: time.Now().UTC(),
	}
	if _, err := doc.Set(ctx, first); err != nil {
		t.Fatalf("writing first: %v", err)
	}

	second := model.Topic{
		Name: "Bruchrechnung",
		Concepts: []model.Concept{
			{ID: "b", Name: "Concept B"},
			{ID: "c", Name: "Concept C", Prerequisites: []string{"b"}},
		},
		AddedAt: time.Now().UTC(),
	}
	if _, err := doc.Set(ctx, second); err != nil {
		t.Fatalf("writing second: %v", err)
	}

	snap, err := doc.Get(ctx)
	if err != nil {
		t.Fatalf("reading: %v", err)
	}
	var got model.Topic
	if err := snap.DataTo(&got); err != nil {
		t.Fatalf("decoding: %v", err)
	}

	if len(got.Concepts) != 2 {
		t.Errorf("expected 2 concepts after overwrite, got %d", len(got.Concepts))
	}
	if got.Concepts[0].ID != "b" {
		t.Errorf("first concept after overwrite: got=%q want=%q", got.Concepts[0].ID, "b")
	}
}

func TestTopicListEmpty(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	uid := "test-topic-empty-user"

	deleteAllTopics(t, client, uid)

	iter := client.Collection("users").Doc(uid).Collection("topics").Documents(ctx)
	snap, err := iter.Next()
	if !errors.Is(err, iterator.Done) {
		if snap != nil {
			t.Fatalf("expected empty collection, got doc %s (err=%v)", snap.Ref.ID, err)
		}
		t.Fatalf("expected iterator.Done, got err=%v", err)
	}
}
