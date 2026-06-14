package mission

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	options "github.com/geoberle/mathilde/backend/cmd"
)

func DefaultOptions() *RawMissionOptions {
	return &RawMissionOptions{
		Options: options.DefaultOptions(),
	}
}

func BindOptions(opts *RawMissionOptions, cmd *cobra.Command) {
	options.BindOptions(opts.Options, cmd)
}

type RawMissionOptions struct {
	Options *options.RawOptions
}

func (o *RawMissionOptions) Validate() (*ValidatedMissionOptions, error) {
	validated, err := o.Options.Validate()
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return &ValidatedMissionOptions{
		validatedMissionOptions: &validatedMissionOptions{
			ValidatedOptions: validated,
		},
	}, nil
}

type validatedMissionOptions struct {
	*options.ValidatedOptions
}

type ValidatedMissionOptions struct {
	*validatedMissionOptions
}

func (o *ValidatedMissionOptions) Complete(ctx context.Context) (*MissionOptions, error) {
	completed, err := o.ValidatedOptions.Complete(ctx)
	if err != nil {
		return nil, err
	}
	return &MissionOptions{
		completedMissionOptions: &completedMissionOptions{
			Options: completed,
		},
	}, nil
}

type completedMissionOptions struct {
	*options.Options
}

type MissionOptions struct {
	*completedMissionOptions
}
