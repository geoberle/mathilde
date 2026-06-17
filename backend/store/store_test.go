package store_test

import (
	"context"
	"os"
	"testing"

	"cloud.google.com/go/firestore"

	"github.com/geoberle/mathilde/backend/store"
)

func newTestStore(t *testing.T) *store.Store {
	t.Helper()

	emulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST")
	if len(emulatorHost) == 0 {
		t.Skip("FIRESTORE_EMULATOR_HOST not set — skipping integration test (run: gcloud emulators firestore start)")
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, store.ProjectID())
	if err != nil {
		t.Fatalf("creating firestore client: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	return &store.Store{Client: client}
}

func newPathTestStore(t *testing.T) *store.Store {
	t.Helper()
	ctx := context.Background()

	t.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:1")

	client, err := firestore.NewClient(ctx, store.ProjectID())
	if err != nil {
		t.Fatalf("creating firestore client: %v", err)
	}
	t.Cleanup(func() { client.Close() })

	return &store.Store{Client: client}
}

func TestRoundTrip(t *testing.T) {
	s := newTestStore(t)
	ctx := context.Background()
	uid := "test-user"

	doc := s.Doc(uid, "ping", "test")

	_, err := doc.Set(ctx, map[string]any{
		"message": "hello from test",
	})
	if err != nil {
		t.Fatalf("writing document: %v", err)
	}

	snap, err := doc.Get(ctx)
	if err != nil {
		t.Fatalf("reading document: %v", err)
	}

	data := snap.Data()
	got, _ := data["message"].(string)
	if got != "hello from test" {
		t.Errorf("got message=%q, want %q", got, "hello from test")
	}

	_, err = doc.Delete(ctx)
	if err != nil {
		t.Fatalf("deleting document: %v", err)
	}
}

func TestCollectionPaths(t *testing.T) {
	s := newPathTestStore(t)
	uid := "test-user"

	tests := []struct {
		name       string
		collection string
	}{
		{name: "profile", collection: "profile"},
		{name: "topics", collection: "topics"},
		{name: "sessions", collection: "sessions"},
		{name: "results", collection: "results"},
		{name: "learning-records", collection: "learning-records"},
		{name: "evaluations", collection: "evaluations"},
		{name: "reference", collection: "reference"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			col := s.Collection(uid, tt.collection)
			wantSuffix := "users/" + uid + "/" + tt.collection
			if got := col.Path; len(got) < len(wantSuffix) || got[len(got)-len(wantSuffix):] != wantSuffix {
				t.Errorf("got path=%q, want suffix %q", got, wantSuffix)
			}
		})
	}
}

func TestProfileDoc(t *testing.T) {
	s := newPathTestStore(t)
	uid := "test-user"

	doc := s.ProfileDoc(uid)
	wantSuffix := "users/" + uid + "/profile/main"
	if got := doc.Path; len(got) < len(wantSuffix) || got[len(got)-len(wantSuffix):] != wantSuffix {
		t.Errorf("got path=%q, want suffix %q", got, wantSuffix)
	}
}

func TestProgressDoc(t *testing.T) {
	s := newPathTestStore(t)
	uid := "test-user"

	doc := s.ProgressDoc(uid)
	wantSuffix := "users/" + uid + "/progress/main"
	if got := doc.Path; len(got) < len(wantSuffix) || got[len(got)-len(wantSuffix):] != wantSuffix {
		t.Errorf("got path=%q, want suffix %q", got, wantSuffix)
	}
}
