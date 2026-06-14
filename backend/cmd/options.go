package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/geoberle/mathilde/backend/store"
)

func DefaultOptions() *RawOptions {
	return &RawOptions{}
}

func BindOptions(opts *RawOptions, cmd *cobra.Command) {
	cmd.Flags().StringVar(&opts.UID, "uid", opts.UID, "Firebase user UID")
}

// RawOptions holds shared CLI flag values.
type RawOptions struct {
	UID string
}

func (o *RawOptions) Validate() (*ValidatedOptions, error) {
	if len(o.UID) == 0 {
		return nil, fmt.Errorf("--uid is required")
	}
	return &ValidatedOptions{
		validatedOptions: &validatedOptions{
			RawOptions: o,
		},
	}, nil
}

type validatedOptions struct {
	*RawOptions
}

type ValidatedOptions struct {
	*validatedOptions
}

func (o *ValidatedOptions) Complete(ctx context.Context) (*Options, error) {
	s, err := store.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("connecting to Firestore: %w", err)
	}
	return &Options{
		completedOptions: &completedOptions{
			Store: s,
			UID:   o.UID,
		},
	}, nil
}

type completedOptions struct {
	Store *store.Store
	UID   string
}

type Options struct {
	*completedOptions
}

func (o *Options) Close() error {
	return o.Store.Close()
}
