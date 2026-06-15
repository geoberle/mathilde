package store

import (
	"context"
	"strings"

	"github.com/geoberle/mathilde/backend/model"
)

const topicsCollection = "topics"

func (s *Store) GetTopic(ctx context.Context, uid, name string) (*Document[model.Topic], error) {
	return Get[model.Topic](ctx, s.Doc(uid, topicsCollection, slugify(name)))
}

func (s *Store) ListTopics(ctx context.Context, uid string) ([]Document[model.Topic], error) {
	return List[model.Topic](ctx, s.Collection(uid, topicsCollection))
}

func (s *Store) CreateTopic(ctx context.Context, uid, name string, topic *model.Topic) (*Document[model.Topic], error) {
	return Create(ctx, s.Doc(uid, topicsCollection, slugify(name)), topic)
}

func (s *Store) ReplaceTopic(ctx context.Context, uid, name string, doc *Document[model.Topic]) (*Document[model.Topic], error) {
	return Replace(ctx, s.Doc(uid, topicsCollection, slugify(name)), doc)
}

// TopicSlug returns the document ID for a topic name.
func TopicSlug(name string) string {
	return slugify(name)
}

func slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "ü", "ue")
	s = strings.ReplaceAll(s, "ß", "ss")
	return s
}
