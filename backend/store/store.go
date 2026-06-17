package store

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
)

const DefaultProjectID = "mathilde-61d77"

func ProjectID() string {
	if id := os.Getenv("FIRESTORE_PROJECT_ID"); len(id) > 0 {
		return id
	}
	return DefaultProjectID
}

// Store wraps a Firestore client scoped to a specific user.
type Store struct {
	Client *firestore.Client
}

// New creates a Store using Firebase Admin SDK.
// Credentials come from GOOGLE_APPLICATION_CREDENTIALS env var.
func New(ctx context.Context) (*Store, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: ProjectID()})
	if err != nil {
		return nil, fmt.Errorf("initializing firebase app: %w", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating firestore client: %w", err)
	}

	return &Store{Client: client}, nil
}

// Close shuts down the Firestore client.
func (s *Store) Close() error {
	return s.Client.Close()
}

// UserDoc returns the base document reference for a user.
func (s *Store) UserDoc(uid string) *firestore.DocumentRef {
	return s.Client.Collection("users").Doc(uid)
}

// ProfileDoc returns the profile document reference for a user.
func (s *Store) ProfileDoc(uid string) *firestore.DocumentRef {
	return s.UserDoc(uid).Collection("profile").Doc("main")
}

// ProgressDoc returns the progress document reference for a user.
func (s *Store) ProgressDoc(uid string) *firestore.DocumentRef {
	return s.UserDoc(uid).Collection("progress").Doc("main")
}

// Collection returns a collection reference under a user.
func (s *Store) Collection(uid, name string) *firestore.CollectionRef {
	return s.UserDoc(uid).Collection(name)
}

// Doc returns a document reference within a user's collection.
func (s *Store) Doc(uid, collection, docID string) *firestore.DocumentRef {
	return s.Collection(uid, collection).Doc(docID)
}
