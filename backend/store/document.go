package store

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound = errors.New("document not found")
	ErrConflict = errors.New("document was modified since last read")
)

// Document wraps a Firestore document's data with its update timestamp,
// which serves as a precondition for conditional Replace operations.
type Document[T any] struct {
	ID         string
	Data       T
	UpdateTime time.Time
}

// Get retrieves a single document by reference.
// Returns ErrNotFound if the document does not exist.
func Get[T any](ctx context.Context, ref *firestore.DocumentRef) (*Document[T], error) {
	snap, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("reading %s: %w", ref.Path, err)
	}
	var data T
	if err := snap.DataTo(&data); err != nil {
		return nil, fmt.Errorf("decoding %s: %w", ref.Path, err)
	}
	return &Document[T]{
		ID:         snap.Ref.ID,
		Data:       data,
		UpdateTime: snap.UpdateTime,
	}, nil
}

// List returns all documents in a collection.
func List[T any](ctx context.Context, col *firestore.CollectionRef) ([]Document[T], error) {
	snaps, err := col.Documents(ctx).GetAll()
	if err != nil {
		return nil, fmt.Errorf("listing %s: %w", col.Path, err)
	}
	docs := make([]Document[T], 0, len(snaps))
	for _, snap := range snaps {
		var data T
		if err := snap.DataTo(&data); err != nil {
			return nil, fmt.Errorf("decoding %s: %w", snap.Ref.Path, err)
		}
		docs = append(docs, Document[T]{
			ID:         snap.Ref.ID,
			Data:       data,
			UpdateTime: snap.UpdateTime,
		})
	}
	return docs, nil
}

// Create writes a new document. Fails if the document already exists.
func Create[T any](ctx context.Context, ref *firestore.DocumentRef, data *T) (*Document[T], error) {
	wr, err := ref.Create(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("creating %s: %w", ref.Path, err)
	}
	return &Document[T]{
		ID:         ref.ID,
		Data:       *data,
		UpdateTime: wr.UpdateTime,
	}, nil
}

// Replace overwrites an existing document, conditional on its UpdateTime
// matching the value from a prior Get. Returns the updated document.
func Replace[T any](ctx context.Context, ref *firestore.DocumentRef, doc *Document[T]) (*Document[T], error) {
	updates := structToUpdates(doc.Data)
	wr, err := ref.Update(ctx, updates, firestore.LastUpdateTime(doc.UpdateTime))
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, fmt.Errorf("replacing %s: %w", ref.Path, ErrNotFound)
		}
		if status.Code(err) == codes.FailedPrecondition {
			return nil, fmt.Errorf("replacing %s: %w", ref.Path, ErrConflict)
		}
		return nil, fmt.Errorf("replacing %s: %w", ref.Path, err)
	}
	return &Document[T]{
		ID:         doc.ID,
		Data:       doc.Data,
		UpdateTime: wr.UpdateTime,
	}, nil
}

// structToUpdates converts a struct with firestore tags to []firestore.Update.
// Fields tagged with omitempty that have zero values are deleted from the document.
func structToUpdates(v any) []firestore.Update {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	t := val.Type()
	var updates []firestore.Update
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("firestore")
		if len(tag) == 0 || tag == "-" {
			continue
		}
		name, opts, _ := strings.Cut(tag, ",")
		fieldVal := val.Field(i)
		if strings.Contains(opts, "omitempty") && fieldVal.IsZero() {
			updates = append(updates, firestore.Update{
				Path:  name,
				Value: firestore.Delete,
			})
			continue
		}
		updates = append(updates, firestore.Update{
			Path:  name,
			Value: fieldVal.Interface(),
		})
	}
	return updates
}
