package topic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/spf13/cobra"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/geoberle/mathilde/backend/model"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "topic",
		Short: "Manage topics and concepts",
	}

	addCmd, err := newAddCommand()
	if err != nil {
		return nil, err
	}
	listCmd, err := newListCommand()
	if err != nil {
		return nil, err
	}
	showCmd, err := newShowCommand()
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(addCmd, listCmd, showCmd)
	return cmd, nil
}

func newAddCommand() (*cobra.Command, error) {
	opts := DefaultOptions()
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add a new topic with AI concept breakdown",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(cmd.Context(), opts, args[0])
		},
	}
	opts.BindOptions(cmd)
	return cmd, nil
}

func newListCommand() (*cobra.Command, error) {
	opts := DefaultOptions()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all topics with concept count",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd.Context(), opts)
		},
	}
	opts.BindOptions(cmd)
	return cmd, nil
}

func newShowCommand() (*cobra.Command, error) {
	opts := DefaultOptions()
	cmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Show concept tree for a topic",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShow(cmd.Context(), opts, args[0])
		},
	}
	opts.BindOptions(cmd)
	return cmd, nil
}

func runAdd(ctx context.Context, opts *RawTopicOptions, name string) error {
	validated, err := opts.Validate()
	if err != nil {
		return err
	}
	completed, err := validated.Complete(ctx)
	if err != nil {
		return err
	}
	defer completed.Close()

	profile, err := readProfile(ctx, completed)
	if err != nil {
		return fmt.Errorf("reading profile: %w", err)
	}

	fmt.Printf("Generating concept breakdown for %q...\n", name)

	concepts, err := generateConcepts(ctx, completed.Claude, name, profile)
	if err != nil {
		return fmt.Errorf("generating concepts: %w", err)
	}

	topic := model.Topic{
		Name:     name,
		Concepts: concepts,
		AddedAt:  time.Now().UTC(),
	}

	docID := slugify(name)
	doc := completed.Store.Doc(completed.UID, "topics", docID)
	if _, err := doc.Set(ctx, topic); err != nil {
		return fmt.Errorf("writing topic: %w", err)
	}

	fmt.Printf("Topic %q added with %d concepts:\n", name, len(concepts))
	for _, c := range concepts {
		prereqs := ""
		if len(c.Prerequisites) > 0 {
			prereqs = fmt.Sprintf(" (requires: %s)", strings.Join(c.Prerequisites, ", "))
		}
		fmt.Printf("  - %s%s\n", c.Name, prereqs)
	}
	return nil
}

func runList(ctx context.Context, opts *RawTopicOptions) error {
	validated, err := opts.Validate()
	if err != nil {
		return err
	}
	completed, err := validated.Complete(ctx)
	if err != nil {
		return err
	}
	defer completed.Close()

	iter := completed.Store.Collection(completed.UID, "topics").Documents(ctx)
	found := false
	for {
		snap, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return fmt.Errorf("reading topics: %w", err)
		}
		var t model.Topic
		if err := snap.DataTo(&t); err != nil {
			return fmt.Errorf("decoding topic %s: %w", snap.Ref.ID, err)
		}
		fmt.Printf("%-30s %d concepts\n", t.Name, len(t.Concepts))
		found = true
	}
	if !found {
		fmt.Println("No topics yet. Use 'mathilde topic add <name>' to add one.")
	}
	return nil
}

func runShow(ctx context.Context, opts *RawTopicOptions, name string) error {
	validated, err := opts.Validate()
	if err != nil {
		return err
	}
	completed, err := validated.Complete(ctx)
	if err != nil {
		return err
	}
	defer completed.Close()

	topic, err := findTopicByName(ctx, completed, name)
	if err != nil {
		return err
	}

	printTopicTree(os.Stdout, topic)
	return nil
}

func printTopicTree(w io.Writer, topic *model.Topic) {
	fmt.Fprintf(w, "Topic: %s (%d concepts)\n\n", topic.Name, len(topic.Concepts))

	conceptNames := make(map[string]string)
	for _, c := range topic.Concepts {
		conceptNames[c.ID] = c.Name
	}

	for _, c := range topic.Concepts {
		fmt.Fprintf(w, "  %s\n", c.Name)
		fmt.Fprintf(w, "    ID: %s\n", c.ID)
		if len(c.Prerequisites) > 0 {
			prereqNames := make([]string, 0, len(c.Prerequisites))
			for _, pid := range c.Prerequisites {
				if pname, ok := conceptNames[pid]; ok {
					prereqNames = append(prereqNames, pname)
				} else {
					prereqNames = append(prereqNames, pid)
				}
			}
			fmt.Fprintf(w, "    Requires: %s\n", strings.Join(prereqNames, ", "))
		}
		fmt.Fprintln(w)
	}
}

func readProfile(ctx context.Context, opts *TopicOptions) (*model.Profile, error) {
	snap, err := opts.Store.ProfileDoc(opts.UID).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return &model.Profile{}, nil
		}
		return nil, fmt.Errorf("reading profile: %w", err)
	}
	var p model.Profile
	if err := snap.DataTo(&p); err != nil {
		return nil, fmt.Errorf("decoding profile: %w", err)
	}
	return &p, nil
}

func findTopicByName(ctx context.Context, opts *TopicOptions, name string) (*model.Topic, error) {
	slug := slugify(name)
	doc := opts.Store.Doc(opts.UID, "topics", slug)
	snap, err := doc.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("topic %q not found (looked up %q): %w", name, slug, err)
	}
	var t model.Topic
	if err := snap.DataTo(&t); err != nil {
		return nil, fmt.Errorf("decoding topic: %w", err)
	}
	return &t, nil
}

func slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "ä", "ae")
	s = strings.ReplaceAll(s, "ö", "oe")
	s = strings.ReplaceAll(s, "ü", "ue")
	s = strings.ReplaceAll(s, "ß", "ss")
	return s
}

const conceptPrompt = `Du bist ein erfahrener österreichischer Mathematiklehrer am Gymnasium.

Zerlege das Thema "%[1]s" in seine atomaren Lernkonzepte. Jedes Konzept soll eine einzelne Lernsitzung (10-15 Minuten) abdecken.

Kontext:
- Schulstufe: Gymnasium, österreichischer Lehrplan
%[2]s
Regeln:
- Verwende österreichische mathematische Fachbegriffe (z.B. "Bruch" nicht "Fraktion")
- Ordne die Konzepte so, dass Voraussetzungen zuerst kommen
- Jedes Konzept hat eine eindeutige ID (kebab-case, deutsch)
- prerequisites enthält die IDs der Konzepte, die vorher verstanden sein müssen
- Gib NUR das JSON-Array zurück, keinen anderen Text

Beispiel-Format:
[
  {"id": "was-ist-ein-bruch", "name": "Was ist ein Bruch?", "prerequisites": []},
  {"id": "erweitern", "name": "Erweitern von Brüchen", "prerequisites": ["was-ist-ein-bruch"]}
]

Thema: %[1]s`

func buildProfileContext(profile *model.Profile) string {
	var parts []string
	if grade, ok := profile.Preferences["grade"]; ok {
		parts = append(parts, fmt.Sprintf("- Schulstufe: %s", grade))
	}
	if len(profile.Mission) > 0 {
		parts = append(parts, fmt.Sprintf("- Lernziel: %s", profile.Mission))
	}
	return strings.Join(parts, "\n")
}

func generateConcepts(ctx context.Context, client *anthropic.Client, topicName string, profile *model.Profile) ([]model.Concept, error) {
	profileCtx := buildProfileContext(profile)
	prompt := fmt.Sprintf(conceptPrompt, topicName, profileCtx)

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("calling Claude API: %w", err)
	}

	if len(message.Content) == 0 {
		return nil, fmt.Errorf("empty response from AI")
	}

	return parseConceptsFromText(message.Content[0].Text)
}

func parseConceptsFromText(text string) ([]model.Concept, error) {
	text = strings.TrimSpace(text)

	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		start := 1
		end := len(lines) - 1
		if end > start {
			text = strings.Join(lines[start:end], "\n")
		}
	}

	var concepts []model.Concept
	if err := json.Unmarshal([]byte(text), &concepts); err != nil {
		return nil, fmt.Errorf("parsing concept JSON: %w\nraw text:\n%s", err, text)
	}

	if len(concepts) == 0 {
		return nil, fmt.Errorf("no concepts in response")
	}

	for i := range concepts {
		if len(concepts[i].Status) == 0 {
			concepts[i].Status = model.StatusNotStarted
		}
	}

	return concepts, nil
}
