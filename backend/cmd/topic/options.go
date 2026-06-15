package topic

import (
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/vertex"
	"github.com/spf13/cobra"

	options "github.com/geoberle/mathilde/backend/cmd"
	"github.com/geoberle/mathilde/backend/store"
)

func DefaultOptions() *RawTopicOptions {
	return &RawTopicOptions{
		Options:       options.DefaultOptions(),
		VertexRegion:  vertexRegion(),
		VertexProject: vertexProjectID(),
	}
}

func (o *RawTopicOptions) BindOptions(cmd *cobra.Command) {
	options.BindOptions(o.Options, cmd)
}

type RawTopicOptions struct {
	Options       *options.RawOptions
	VertexRegion  string
	VertexProject string
}

func (o *RawTopicOptions) Validate() (*ValidatedTopicOptions, error) {
	validated, err := o.Options.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return &ValidatedTopicOptions{
		validatedTopicOptions: &validatedTopicOptions{
			ValidatedOptions: validated,
			vertexRegion:     o.VertexRegion,
			vertexProject:    o.VertexProject,
		},
	}, nil
}

type validatedTopicOptions struct {
	*options.ValidatedOptions
	vertexRegion  string
	vertexProject string
}

type ValidatedTopicOptions struct {
	*validatedTopicOptions
}

func vertexRegion() string {
	if r := os.Getenv("VERTEX_REGION"); len(r) > 0 {
		return r
	}
	return "europe-west1"
}

func vertexProjectID() string {
	if id := os.Getenv("VERTEX_PROJECT_ID"); len(id) > 0 {
		return id
	}
	return store.ProjectID()
}

func (o *ValidatedTopicOptions) Complete(ctx context.Context) (*TopicOptions, error) {
	completed, err := o.ValidatedOptions.Complete(ctx)
	if err != nil {
		return nil, err
	}
	claude := anthropic.NewClient(
		vertex.WithGoogleAuth(ctx, o.vertexRegion, o.vertexProject,
			"https://www.googleapis.com/auth/cloud-platform",
		),
	)
	return &TopicOptions{
		completedTopicOptions: &completedTopicOptions{
			Options: completed,
			Claude:  &claude,
		},
	}, nil
}

type completedTopicOptions struct {
	*options.Options
	Claude *anthropic.Client
}

type TopicOptions struct {
	*completedTopicOptions
}
