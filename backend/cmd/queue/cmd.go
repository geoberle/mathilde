package queue

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "queue",
		Short: "Manage the session queue",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List queued sessions",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented — see issue #12")
		},
	})

	return cmd, nil
}
