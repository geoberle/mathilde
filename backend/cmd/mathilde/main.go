package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/geoberle/mathilde/backend/cmd/mission"
	"github.com/geoberle/mathilde/backend/cmd/ping"
	"github.com/geoberle/mathilde/backend/cmd/progress"
	"github.com/geoberle/mathilde/backend/cmd/queue"
	"github.com/geoberle/mathilde/backend/cmd/session"
	"github.com/geoberle/mathilde/backend/cmd/topic"
)

func main() {
	cmd := &cobra.Command{
		Use:           "mathilde",
		Short:         "Mathilde — adaptive math learning CLI for the curator",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	commands := []func() (*cobra.Command, error){
		ping.NewCommand,
		topic.NewCommand,
		session.NewCommand,
		queue.NewCommand,
		progress.NewCommand,
		mission.NewCommand,
	}
	for _, newCmd := range commands {
		c, err := newCmd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		cmd.AddCommand(c)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
