package ping

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/spf13/cobra"
)

func NewCommand() (*cobra.Command, error) {
	opts := DefaultOptions()
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "Test Firestore connectivity with a round-trip read/write",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), opts)
		},
	}
	BindOptions(opts, cmd)
	return cmd, nil
}

func run(ctx context.Context, opts *RawPingOptions) error {
	validated, err := opts.Validate()
	if err != nil {
		return err
	}
	completed, err := validated.Complete(ctx)
	if err != nil {
		return err
	}
	defer completed.Options.Close()
	return completed.Execute(ctx)
}

func (o *PingOptions) Execute(ctx context.Context) error {
	doc := o.Options.Store.Doc(o.Options.UID, "ping", "test")
	now := time.Now().UTC()

	_, err := doc.Set(ctx, map[string]any{
		"message":   "ping from mathilde CLI",
		"timestamp": now,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("writing ping document: %w", err)
	}
	fmt.Println("Write OK")

	snap, err := doc.Get(ctx)
	if err != nil {
		return fmt.Errorf("reading ping document: %w", err)
	}
	data := snap.Data()
	fmt.Printf("Read OK: message=%q timestamp=%v\n", data["message"], data["timestamp"])

	_, err = doc.Delete(ctx)
	if err != nil {
		return fmt.Errorf("deleting ping document: %w", err)
	}
	fmt.Println("Delete OK")
	fmt.Println("Firestore round-trip successful.")
	return nil
}
