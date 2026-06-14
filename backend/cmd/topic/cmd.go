package topic

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "topic",
		Short: "Manage topics and concepts",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "add [name]",
		Short: "Add a new topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented — see issue #4")
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all topics",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented — see issue #4")
		},
	})

	return cmd, nil
}
