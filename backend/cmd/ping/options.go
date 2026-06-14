package ping

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	options "github.com/geoberle/mathilde/backend/cmd"
)

func DefaultOptions() *RawPingOptions {
	return &RawPingOptions{
		Options: options.DefaultOptions(),
	}
}

func BindOptions(opts *RawPingOptions, cmd *cobra.Command) {
	options.BindOptions(opts.Options, cmd)
}

type RawPingOptions struct {
	Options *options.RawOptions
}

func (o *RawPingOptions) Validate() (*ValidatedPingOptions, error) {
	validated, err := o.Options.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return &ValidatedPingOptions{
		validatedPingOptions: &validatedPingOptions{
			ValidatedOptions: validated,
		},
	}, nil
}

type validatedPingOptions struct {
	*options.ValidatedOptions
}

type ValidatedPingOptions struct {
	*validatedPingOptions
}

func (o *ValidatedPingOptions) Complete(ctx context.Context) (*PingOptions, error) {
	completed, err := o.ValidatedOptions.Complete(ctx)
	if err != nil {
		return nil, err
	}
	return &PingOptions{
		completedPingOptions: &completedPingOptions{
			Options: completed,
		},
	}, nil
}

type completedPingOptions struct {
	*options.Options
}

type PingOptions struct {
	*completedPingOptions
}
