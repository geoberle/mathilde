package mission

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/spf13/cobra"

	"github.com/geoberle/mathilde/backend/model"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "mission",
		Short: "View or set the learner's mission",
	}

	getCmd, err := newGetCommand()
	if err != nil {
		return nil, err
	}
	setCmd, err := newSetCommand()
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(getCmd, setCmd)
	return cmd, nil
}

func newGetCommand() (*cobra.Command, error) {
	opts := DefaultOptions()
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Show current mission",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(cmd.Context(), opts)
		},
	}
	BindOptions(opts, cmd)
	return cmd, nil
}

func runGet(ctx context.Context, opts *RawMissionOptions) error {
	validated, err := opts.Validate()
	if err != nil {
		return err
	}
	completed, err := validated.Complete(ctx)
	if err != nil {
		return err
	}
	defer completed.Close()

	snap, err := completed.Store.ProfileDoc(completed.UID).Get(ctx)
	if err != nil {
		return fmt.Errorf("reading profile: %w", err)
	}
	var p model.Profile
	if err := snap.DataTo(&p); err != nil {
		return fmt.Errorf("decoding profile: %w", err)
	}
	if len(p.Mission) == 0 {
		fmt.Println("No mission set.")
		return nil
	}
	fmt.Println(p.Mission)
	return nil
}

func newSetCommand() (*cobra.Command, error) {
	opts := DefaultOptions()
	cmd := &cobra.Command{
		Use:   "set [mission text]",
		Short: "Set the learner's mission",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSet(cmd.Context(), opts, args[0])
		},
	}
	BindOptions(opts, cmd)
	return cmd, nil
}

const maxMissionLength = 500

func runSet(ctx context.Context, opts *RawMissionOptions, text string) error {
	if len(text) == 0 {
		return fmt.Errorf("mission text cannot be empty")
	}
	if len(text) > maxMissionLength {
		return fmt.Errorf("mission text too long (%d chars, max %d)", len(text), maxMissionLength)
	}
	validated, err := opts.Validate()
	if err != nil {
		return err
	}
	completed, err := validated.Complete(ctx)
	if err != nil {
		return err
	}
	defer completed.Close()

	_, err = completed.Store.ProfileDoc(completed.UID).Set(ctx, map[string]any{
		"mission": text,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("setting mission: %w", err)
	}
	fmt.Println("Mission updated.")
	return nil
}
