package store_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"cloud.google.com/go/firestore"

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

func TestGetNotFound(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()

	ref := client.Collection("users").Doc("nonexistent").Collection("profile").Doc("main")
	_, err := store.Get[model.Profile](ctx, ref)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}

func TestCreateAndGet(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	uid := "test-crud-create"

	ref := client.Collection("users").Doc(uid).Collection("profile").Doc("main")
	defer func() { _, _ = ref.Delete(ctx) }()

	want := model.Profile{
		Mission: "Test mission",
	}
	created, err := store.Create(ctx, ref, &want)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.Data.Mission != want.Mission {
		t.Errorf("Create returned mission: got=%q want=%q", created.Data.Mission, want.Mission)
	}
	if created.UpdateTime.IsZero() {
		t.Error("Create returned zero UpdateTime")
	}

	got, err := store.Get[model.Profile](ctx, ref)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Data.Mission != want.Mission {
		t.Errorf("Get mission: got=%q want=%q", got.Data.Mission, want.Mission)
	}
	if got.UpdateTime.IsZero() {
		t.Error("Get returned zero UpdateTime")
	}
}

func TestCreateAlreadyExists(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	uid := "test-crud-exists"

	ref := client.Collection("users").Doc(uid).Collection("profile").Doc("main")
	defer func() { _, _ = ref.Delete(ctx) }()

	profile := model.Profile{Mission: "First"}
	if _, err := store.Create(ctx, ref, &profile); err != nil {
		t.Fatalf("first Create: %v", err)
	}

	second := model.Profile{Mission: "Second"}
	_, err := store.Create(ctx, ref, &second)
	if err == nil {
		t.Fatal("expected error on duplicate Create, got nil")
	}
}

func TestReplaceConditional(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()
	ctx := context.Background()
	uid := "test-crud-replace"

	ref := s.Client.Collection("users").Doc(uid).Collection("profile").Doc("main")
	defer func() { _, _ = ref.Delete(ctx) }()

	original := model.Profile{Mission: "Original"}
	created, err := store.Create(ctx, ref, &original)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	created.Data.Mission = "Updated"
	updated, err := store.Replace(ctx, ref, created)
	if err != nil {
		t.Fatalf("Replace: %v", err)
	}
	if updated.Data.Mission != "Updated" {
		t.Errorf("Replace mission: got=%q want=%q", updated.Data.Mission, "Updated")
	}
	if !updated.UpdateTime.After(created.UpdateTime) {
		t.Error("Replace UpdateTime should be after Create UpdateTime")
	}
}

func TestReplaceConflict(t *testing.T) {
	s := newTestStore(t)
	defer s.Close()
	ctx := context.Background()
	uid := "test-crud-conflict"

	ref := s.Client.Collection("users").Doc(uid).Collection("profile").Doc("main")
	defer func() { _, _ = ref.Delete(ctx) }()

	original := model.Profile{Mission: "Original"}
	doc, err := store.Create(ctx, ref, &original)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Simulate concurrent write by writing directly
	if _, err := ref.Set(ctx, model.Profile{Mission: "Concurrent"}); err != nil {
		t.Fatalf("concurrent write: %v", err)
	}

	doc.Data.Mission = "Stale update"
	_, err = store.Replace(ctx, ref, doc)
	if !errors.Is(err, store.ErrConflict) {
		t.Fatalf("expected ErrConflict, got: %v", err)
	}
}

func TestReplaceNotFound(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()

	ref := client.Collection("users").Doc("test-crud-replace-notfound").Collection("profile").Doc("main")

	doc := &store.Document[model.Profile]{
		ID:         "main",
		Data:       model.Profile{Mission: "Ghost"},
		UpdateTime: time.Now(),
	}

	_, err := store.Replace(ctx, ref, doc)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}

func TestList(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()
	uid := "test-crud-list"

	col := client.Collection("users").Doc(uid).Collection("topics")

	// Clean up
	snaps, _ := col.Documents(ctx).GetAll()
	for _, s := range snaps {
		_, _ = s.Ref.Delete(ctx)
	}
	defer func() {
		snaps, _ := col.Documents(ctx).GetAll()
		for _, s := range snaps {
			_, _ = s.Ref.Delete(ctx)
		}
	}()

	topics := []model.Topic{
		{Name: "Bruchrechnung", Concepts: []model.Concept{{ID: "a", Name: "A"}}, AddedAt: time.Now().UTC()},
		{Name: "Gleichungen", Concepts: []model.Concept{{ID: "b", Name: "B"}}, AddedAt: time.Now().UTC()},
	}
	for _, topic := range topics {
		ref := col.Doc(store.TopicSlug(topic.Name))
		if _, err := ref.Create(ctx, topic); err != nil {
			t.Fatalf("creating topic %q: %v", topic.Name, err)
		}
	}

	docs, err := store.List[model.Topic](ctx, col)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(docs) != 2 {
		t.Fatalf("List count: got=%d want=2", len(docs))
	}
	for _, doc := range docs {
		if doc.UpdateTime.IsZero() {
			t.Errorf("topic %q has zero UpdateTime", doc.Data.Name)
		}
	}
}

func TestListEmpty(t *testing.T) {
	client := newTestClient(t)
	defer client.Close()
	ctx := context.Background()

	col := client.Collection("users").Doc("test-crud-empty").Collection("topics")
	docs, err := store.List[model.Topic](ctx, col)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(docs) != 0 {
		t.Errorf("expected empty list, got %d", len(docs))
	}
}
