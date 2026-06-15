package store

import (
	"context"

	"github.com/geoberle/mathilde/backend/model"
)

func (s *Store) GetProfile(ctx context.Context, uid string) (*Document[model.Profile], error) {
	return Get[model.Profile](ctx, s.ProfileDoc(uid))
}

func (s *Store) CreateProfile(ctx context.Context, uid string, profile *model.Profile) (*Document[model.Profile], error) {
	return Create(ctx, s.ProfileDoc(uid), profile)
}

func (s *Store) ReplaceProfile(ctx context.Context, uid string, doc *Document[model.Profile]) (*Document[model.Profile], error) {
	return Replace(ctx, s.ProfileDoc(uid), doc)
}
