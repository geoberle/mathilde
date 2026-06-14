package session

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Manage sessions",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "review",
		Short: "Review pending sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented — see issue #6")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "approve [session-id]",
		Short: "Approve a session for the learner",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented — see issue #6")
		},
	})

	return cmd, nil
}
