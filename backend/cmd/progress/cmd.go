package progress

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewCommand() (*cobra.Command, error) {
	return &cobra.Command{
		Use:   "progress",
		Short: "Show learner progress",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not implemented — requires learning records")
		},
	}, nil
}
